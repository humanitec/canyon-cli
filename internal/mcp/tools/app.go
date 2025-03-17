package tools

import (
	"context"
	"net/http"

	"github.com/humanitec/humanitec-go-autogen/client"

	"github.com/humanitec/canyon-cli/internal/clients/humanitec"
	"github.com/humanitec/canyon-cli/internal/mcp"
)

func NewGetHumanitecDeploymentSets() mcp.Tool {
	return mcp.Tool{
		Name:        "get_humanitec_deployment_sets",
		Description: `This tool returns the contents of the specified Humanitec Deployment Sets. This can be used to fetch multiple Deployment Sets at once.`,
		InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{
			"org_id":  map[string]interface{}{"type": "string", "description": "The Humanitec Organization (org) ID to work with."},
			"app_id":  map[string]interface{}{"type": "string", "description": "The Humanitec Application (app) ID to work with."},
			"set_ids": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "The list of Humanitec Deployment Set (set) IDs to fetch the contents for."},
		}, "required": []string{"org_id", "app_id", "set_ids"}},
		Callable: func(ctx context.Context, arguments map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			orgId := arguments["org_id"].(string)
			appId := arguments["app_id"].(string)
			setIds := arguments["set_ids"].([]interface{})
			hc, err := humanitec.NewHumanitecClientWithCurrentToken(ctx)
			if err != nil {
				return nil, err
			}
			output := make([]mcp.CallToolResponseContent, 0)
			for _, i := range setIds {
				setId := i.(string)
				if r, err := humanitec.CheckResponse(func() (*client.GetSetResponse, error) {
					return hc.GetSetWithResponse(ctx, orgId, appId, setId, &client.GetSetParams{})
				}).AndStatusCodeEq(http.StatusOK).RespAndError(); err != nil {
					output = append(output, mcp.NewTextToolResponseContent("Failed to fetch contents for set %s: %v", setId, err.Error()))
				} else {
					output = append(output, mcp.NewTextToolResponseContent("The contents of set %s in JSON is: %s", setId, r.Body))
				}
			}
			return output, nil
		},
	}
}
