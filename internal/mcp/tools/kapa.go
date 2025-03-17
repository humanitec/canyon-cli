package tools

import (
	"context"
	"net/http"

	"github.com/humanitec/canyon-cli/internal/clients/humanitec"
	"github.com/humanitec/canyon-cli/internal/mcp"
)

func NewKapaAiDocsTool() mcp.Tool {
	return mcp.Tool{
		Name:        "query_humanitec_documentation",
		Description: `This tool provides access to an LLM that has been fine tuned on Humanitec Platform Orchestrator documentation. This tool provides access to an expert in Humanitec platform engineer. Use this tool whenever you are unsure, need more up to date documentation, or hallucination is a risk.`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{"type": "string"},
			},
			"required": []string{"query"},
		},
		Callable: func(ctx context.Context, m map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			hc, err := humanitec.NewHumanitecClientWithCurrentToken(ctx)
			if err != nil {
				return nil, err
			}
			if r, err := humanitec.CheckResponse(func() (*humanitec.QueryAiDocsResponse, error) {
				return hc.QueryAiDocs(ctx, m["query"].(string))
			}).AndStatusCodeEq(http.StatusOK).RespAndError(); err != nil {
				return nil, err
			} else {
				return []mcp.CallToolResponseContent{
					mcp.NewTextToolResponseContent(r.JSON200.Answer, "assistant"),
				}, nil
			}
		},
	}
}
