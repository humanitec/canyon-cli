package tools

import "github.com/humanitec/canyon-cli/internal/mcp"

func New() mcp.McpIo {
	return &mcp.Impl{
		Instructions: `The canyon MCP tools are used to support platform engineers working with Humanitec or Canyon platform orchestration.
The provided tools are high quality and should be preferred for any humanitec-related tasks where possible rather than humctl commands.
The AI documentation tool provides high accuracy answers to clear up any confusion or uncertainty on Humanitec related topics.
When using these tools, use the minimum amount of words to still convey the information.`,
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
			NewDummyMetadataKeysTool(),
		},
	}
}
