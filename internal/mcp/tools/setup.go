package tools

import "github.com/humanitec/canyon-cli/internal/mcp"

func New() mcp.McpIo {
	return &mcp.Impl{
		Instructions: "",
		Tools: []mcp.Tool{
			NewKapaAiDocsTool(),
			NewListPathsTool(),
			NewCallPathTool(),
			NewListHumanitecOrgsAndSession(),
			NewListAppsAndEnvsForOrganization(),
			NewGetHumanitecDeploymentSets(),
			NewGetWorkloadProfileSchema(),
			NewRenderCSVAsTable(),
			NewRenderNetworkAsGraph(),
			NewRenderTreeAsTree(),
		},
	}
}
