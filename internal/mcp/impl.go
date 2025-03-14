package mcp

import (
	"context"
	"path/filepath"
	"runtime/debug"

	"github.com/humanitec/canyon-cli/internal/rpc"
)

type Impl struct {
	Instructions string
}

var _ McpIo = (*Impl)(nil)

func (m *Impl) SetNotifications(notifications chan<- ServerNotification) {

}

func (m *Impl) SetLevel(ctx context.Context, request SetLevelRequest) (*SetLevelResponse, error) {
	return &SetLevelResponse{}, nil
}

func (m *Impl) GetPrompt(ctx context.Context, request GetPromptRequest) (*GetPromptResponse, error) {
	return nil, rpc.JsonRpcError{Code: -32602, Message: "Unknown prompt"}
}

func (m *Impl) ReadResource(ctx context.Context, request ReadResourceRequest) (*ReadResourceResponse, error) {
	return nil, rpc.JsonRpcError{Code: -32002, Message: "Unknown resource"}
}

func (m *Impl) ListResourcesTemplates(ctx context.Context, request ListResourceTemplatesRequest) (*ListResourceTemplatesResponse, error) {
	return &ListResourceTemplatesResponse{Resources: []ResourceTemplate{}}, nil
}

func (m *Impl) ListPrompts(ctx context.Context, request ListPromptsRequest) (*ListPromptsResponse, error) {
	return &ListPromptsResponse{Prompts: []Prompt{}}, nil
}

func (m *Impl) ListResources(ctx context.Context, request ListResourcesRequest) (*ListResourcesResponse, error) {
	return &ListResourcesResponse{Resources: []Resource{}}, nil
}

func (m *Impl) Initialize(ctx context.Context, request InitializeRequest) (*InitializeResponse, error) {
	bi, _ := debug.ReadBuildInfo()
	return &InitializeResponse{
		ProtocolVersion: request.ProtocolVersion,
		ServerInfo:      Implementation{Name: filepath.Base(bi.Main.Path), Version: bi.Main.Version},
		Instructions:    m.Instructions,
		Capabilities: ServerCapabilities{
			Tools:     ServerToolsCapabilities{},
			Resources: ServerResourcesCapabilities{},
			Prompts:   ServerPromptsCapabilities{},
		},
	}, nil
}

func (m *Impl) ListTools(ctx context.Context, request ListToolsRequest) (*ListToolsResponse, error) {
	return &ListToolsResponse{
		Tools: []Tool{{
			Name:        "thing",
			Description: "Do thing",
			InputSchema: map[string]interface{}{"type": "object"},
		}},
	}, nil
}

func (m *Impl) CallTool(ctx context.Context, request CallToolRequest) (*CallToolResponse, error) {
	return nil, rpc.JsonRpcError{Code: rpc.JsonRpcInvalidRequestError, Message: "tool not found"}
}
