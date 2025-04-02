package tools

import (
	"context"

	"github.com/humanitec/canyon-cli/internal"
	"github.com/humanitec/canyon-cli/internal/mcp"
)

func NewDummyMetadataKeysTool() mcp.Tool {
	return mcp.Tool{
		Name:        "list_organization_metadata_keys",
		Description: `This tool lists the known metadata keys for an organization. The metadata values for workloads are found in the contents of the score spec, in the deployment set, or on resources.`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"org_id": map[string]interface{}{"type": "string", "description": "The Humanitec Organization (org) ID to work with."},
			},
			"required": []string{"org_id"},
		},
		Callable: func(ctx context.Context, m map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			type fake struct {
				Key         string `json:"key"`
				Description string `json:"description"`
			}
			raw := internal.PrettyJson([]fake{
				{Key: "Service-Owner", Description: "The project team who own this workload and are responsible for development and deployments"},
				{Key: "Github-Repo-Url", Description: "The GitHub repository URL where the source code of a workload can be found"},
				{Key: "Git-Tag", Description: "The Git tag that the container image of this workload comes from"},
				{Key: "Grafana-Dashboard-Url", Description: "The grafana dashboard for the operations metrics"},
				{Key: "Aws-Arn", Description: "The AWS ARN id of the related resource"},
			})
			return []mcp.CallToolResponseContent{
				mcp.NewTextToolResponseContent("The following workload and resource metadata keys are known for this org in JSON format: %s", string(raw)),
			}, nil
		},
	}
}
