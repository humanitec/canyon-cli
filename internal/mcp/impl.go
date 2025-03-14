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

func (m *Impl) SetLevel(ctx context.Context, request SetLevelRequest, notifications chan<- ServerNotification) (*SetLevelResponse, error) {
	return &SetLevelResponse{}, nil
}

func (m *Impl) GetPrompt(ctx context.Context, request GetPromptRequest, notifications chan<- ServerNotification) (*GetPromptResponse, error) {
	return nil, rpc.JsonRpcError{Code: -32602, Message: "Unknown prompt"}
}

func (m *Impl) ReadResource(ctx context.Context, request ReadResourceRequest, notifications chan<- ServerNotification) (*ReadResourceResponse, error) {
	return nil, rpc.JsonRpcError{Code: -32002, Message: "Unknown resource"}
}

func (m *Impl) ListResourcesTemplates(ctx context.Context, request ListResourceTemplatesRequest, notifications chan<- ServerNotification) (*ListResourceTemplatesResponse, error) {
	return &ListResourceTemplatesResponse{Resources: []ResourceTemplate{}}, nil
}

func (m *Impl) ListPrompts(ctx context.Context, request ListPromptsRequest, notifications chan<- ServerNotification) (*ListPromptsResponse, error) {
	return &ListPromptsResponse{Prompts: []Prompt{}}, nil
}

func (m *Impl) ListResources(ctx context.Context, request ListResourcesRequest, notifications chan<- ServerNotification) (*ListResourcesResponse, error) {
	return &ListResourcesResponse{Resources: []Resource{}}, nil
}

func (m *Impl) Initialize(ctx context.Context, request InitializeRequest, notifications chan<- ServerNotification) (*InitializeResponse, error) {
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

func (m *Impl) ListTools(ctx context.Context, request ListToolsRequest, notifications chan<- ServerNotification) (*ListToolsResponse, error) {
	return &ListToolsResponse{
		Tools: []Tool{{
			Name:        "thing",
			Description: "Do thing",
			InputSchema: map[string]interface{}{"type": "object"},
		}},
	}, nil
}

func (m *Impl) CallTool(ctx context.Context, request CallToolRequest, notifications chan<- ServerNotification) (*CallToolResponse, error) {
	return nil, rpc.JsonRpcError{Code: rpc.JsonRpcInvalidRequestError, Message: "tool not found"}
}
