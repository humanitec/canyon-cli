package tools

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/humanitec/canyon-cli/internal/clients/humanitec"
	"github.com/humanitec/canyon-cli/internal/mcp"
)

func NewListPathsTool() mcp.Tool {
	return mcp.Tool{
		Name: "list-canyon-paths",
		Description: `Returns a list of 'paths' supported by the canyon MCP server.
Paths are remote functions which can be used to query or achieve a wide array of functionality.
The list of available paths may change over time so consider listing the available paths when there is low confidence that an existing paths can be used to solve the user query.
Canyon paths are not tools themselves and must be called through the call-canyon-path tool.`,
		InputSchema: map[string]interface{}{
			"type":     "object",
			"required": []interface{}{"org_id"},
			"properties": map[string]interface{}{
				"org_id": map[string]interface{}{"type": "string", "description": "The organization ID"},
			}},
		Callable: func(ctx context.Context, arguments map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			hc, err := humanitec.NewHumanitecClientWithCurrentToken(ctx)
			if err != nil {
				return nil, err
			}

			var tools []mcp.ToolResponse
			if sum, err := hc.ListActionPipelineSummaries(ctx, arguments["org_id"].(string)); err != nil {
				return nil, err
			} else if sum.JSON200 == nil {
				// This is a hack for demos while the action pipelines are feature flagged off
				if sum.StatusCode() == http.StatusForbidden || sum.StatusCode() == http.StatusMethodNotAllowed {
					return []mcp.CallToolResponseContent{
						mcp.NewTextToolResponseContent("There are no paths available in this org"),
					}, nil
				}
				return nil, fmt.Errorf("unexpected response from humanitec: %s %s", sum.HTTPResponse.Status, string(sum.Body))
			} else {
				for _, summary := range sum.JSON200 {

					if ap, err := hc.GetActionPipeline(ctx, summary.OrgId, summary.Id); err != nil {
						return nil, err
					} else if ap.JSON200 == nil {
						return nil, fmt.Errorf("unexpected response from humanitec: %v", ap)
					} else {
						tools = append(tools, mcp.ToolResponse{
							Name:        ap.JSON200.Id,
							Description: ap.JSON200.Description,
							InputSchema: ap.JSON200.InputsJsonSchema,
						})
					}
				}
			}

			raw, _ := json.Marshal(tools)
			return []mcp.CallToolResponseContent{
				mcp.NewTextToolResponseContent("Here's an array of the current canyon tools in JSON: %s", string(raw)),
			}, nil
		},
	}
}

func NewCallPathTool() mcp.Tool {
	return mcp.Tool{
		Name:        "call-canyon-path",
		Description: "Call a canyon path previously discovered through list-canyon-paths",
		InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{
			"org_id":          map[string]interface{}{"type": "string", "description": "The organization ID of the org in which the path is defined"},
			"name":            map[string]interface{}{"type": "string", "description": "The name of the path to call"},
			"arguments":       map[string]interface{}{"type": "object", "description": "The arguments of the path to call, these must match the input schema"},
			"idempotency_key": map[string]interface{}{"type": "string", "description": "An idempotency key to use to continue the request if it times out, this will be created for you on the first attempt"},
		}, "required": []interface{}{"org_id", "name", "arguments", "idempotency_key"}},
		Callable: func(ctx context.Context, arguments map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			name, _ := arguments["name"].(string)
			args, _ := arguments["arguments"].(map[string]interface{})

			hc, err := humanitec.NewHumanitecClientWithCurrentToken(ctx)
			if err != nil {
				return nil, err
			}

			idempotencyKey, _ := arguments["idempotency_key"].(string)
			if idempotencyKey == "" {
				idempotencyKeyRaw := make([]byte, 10)
				_, _ = rand.Read(idempotencyKeyRaw)
				idempotencyKey = hex.EncodeToString(idempotencyKeyRaw)
			}

			if r, err := hc.CallActionPipeline(ctx, arguments["org_id"].(string), name, &humanitec.CallActionPipelineParams{IdempotencyKey: idempotencyKey}, humanitec.CallActionPipelineRequestBody{
				Inputs: args,
			}); err != nil {
				return nil, err
			} else if r.JSON200 == nil {
				// This is a hack for demos while the action pipelines are feature flagged off
				if r.StatusCode() == http.StatusForbidden || r.StatusCode() == http.StatusMethodNotAllowed {
					return []mcp.CallToolResponseContent{
						mcp.NewTextToolResponseContent("There are no paths available in this org"),
					}, nil
				}
				if r.StatusCode() == http.StatusGatewayTimeout {
					return []mcp.CallToolResponseContent{
						mcp.NewTextToolResponseContent("The path timed out after executing for some time, you can make the identity request with idempotency key '%s' to continue waiting", idempotencyKey),
					}, nil
				}
				return nil, fmt.Errorf("unexpected response from humanitec, you can be able to continue the request with idempotency key '%s' to continue waiting: %s %s", idempotencyKey, r.HTTPResponse.Status, string(r.Body))
			} else {
				raw, _ := json.Marshal(r.JSON200)
				return []mcp.CallToolResponseContent{
					mcp.NewTextToolResponseContent("The path returned the following result in JSON: %s", string(raw)),
				}, nil
			}
		},
	}
}
