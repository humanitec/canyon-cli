package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/humanitec/canyon-cli/internal/ref"
	"github.com/humanitec/canyon-cli/internal/rpc"
)

type InitializeRequest struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	ClientInfo      Implementation         `json:"clientInfo"`
	Capabilities    map[string]interface{} `json:"capabilities"`
}

type InitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	ServerInfo      Implementation     `json:"serverInfo"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	Instructions    string             `json:"instructions,omitempty"`
}

type Implementation struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	Logging   ServerLoggingCapabilities   `json:"logging"`
	Prompts   ServerPromptsCapabilities   `json:"prompts"`
	Tools     ServerToolsCapabilities     `json:"tools"`
	Resources ServerResourcesCapabilities `json:"resources"`
}

type ServerLoggingCapabilities struct {
}

type ServerPromptsCapabilities struct {
}

type ServerResourcesCapabilities struct {
}

type ServerToolsCapabilities struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// =========================================

type ListToolsRequest struct {
	Cursor string `json:"cursor"`
}

type ListToolsResponse struct {
	NextCursor string `json:"nextCursor,omitempty"`
	Tools      []Tool `json:"tools"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type CallToolRequest struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type CallToolResponse struct {
	IsError  bool                      `json:"is_error,omitempty"`
	Contents []CallToolResponseContent `json:"content"`
}

type CallToolResponseContent struct {
	*TextContent
	*ImageContent
	*EmbeddedResource
}

func (c CallToolResponseContent) MarshalJSON() ([]byte, error) {
	if c.TextContent != nil {
		return json.Marshal(c.TextContent)
	} else if c.ImageContent != nil {
		return json.Marshal(c.ImageContent)
	} else if c.EmbeddedResource != nil {
		return json.Marshal(c.EmbeddedResource)
	} else {
		return []byte("{}"), nil
	}
}

func NewTextToolResponseContent(text string, args ...any) CallToolResponseContent {
	return CallToolResponseContent{TextContent: &TextContent{Text: fmt.Sprintf(text, args...)}}
}

func NewTextToolResponseContentWithAudience(text string, aud string) CallToolResponseContent {
	return CallToolResponseContent{TextContent: &TextContent{Text: text, Annotations: &Annotations{Audience: []string{aud}}}}
}

type TextContentType struct {
}

func (r TextContentType) MarshalJSON() ([]byte, error) {
	return json.Marshal("text")
}

type TextContent struct {
	Type        TextContentType `json:"type"`
	Text        string          `json:"text"`
	Annotations *Annotations    `json:"annotations,omitempty"`
}

type EmbeddedResourceContentType struct {
}

func (t EmbeddedResourceContentType) MarshalJSON() ([]byte, error) {
	return json.Marshal("resource")
}

type EmbeddedResource struct {
	Type        EmbeddedResourceContentType `json:"type"`
	Resource    ResourceContent             `json:"resource"`
	Annotations *Annotations                `json:"annotations,omitempty"`
}

type ImageContentType struct {
}

func (t ImageContentType) MarshalJSON() ([]byte, error) {
	return json.Marshal("image")
}

type ImageContent struct {
	Type        ImageContentType `json:"type"`
	MimeType    string           `json:"mimeType"`
	Data        string           `json:"data"`
	Annotations *Annotations     `json:"annotations,omitempty"`
}

type Annotations struct {
	Audience []string `json:"audience"`
	Priority float64  `json:"priority,omitempty"`
}

// =========================================

type ListPromptsRequest struct {
	Cursor string `json:"cursor"`
}

type ListPromptsResponse struct {
	NextCursor string   `json:"nextCursor,omitempty"`
	Prompts    []Prompt `json:"prompts"`
}

type Prompt struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Arguments   []PromptArgument `json:"arguments,omitempty"`
}

type PromptArgument struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

type GetPromptRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Arguments   map[string]interface{} `json:"arguments"`
}

type GetPromptResponse struct {
	Description string          `json:"description"`
	Messages    []PromptMessage `json:"messages"`
}

type PromptMessage struct {
	Role    string               `json:"role"`
	Content PromptMessageContent `json:"content"`
}

type PromptMessageContent struct {
	*TextContent
	*ImageContent
	*EmbeddedResource
}

func (c PromptMessageContent) MarshalJSON() ([]byte, error) {
	if c.TextContent != nil {
		return json.Marshal(c.TextContent)
	} else if c.ImageContent != nil {
		return json.Marshal(c.ImageContent)
	} else if c.EmbeddedResource != nil {
		return json.Marshal(c.EmbeddedResource)
	} else {
		return []byte("{}"), nil
	}
}

// =========================================

type ListResourcesRequest struct {
	Cursor string `json:"cursor"`
}

type ListResourcesResponse struct {
	NextCursor string     `json:"nextCursor,omitempty"`
	Resources  []Resource `json:"resources"`
}

type Resource struct {
	Uri         string       `json:"uri"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Size        int64        `json:"size,omitempty"`
	MimeType    string       `json:"mimeType,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

type ResourceTemplate struct {
	Name        string       `json:"name"`
	UriTemplate string       `json:"uriTemplate"`
	Description string       `json:"description,omitempty"`
	MimeType    string       `json:"mimeType,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

type ListResourceTemplatesRequest struct {
	Cursor string `json:"cursor"`
}

type ListResourceTemplatesResponse struct {
	NextCursor string             `json:"nextCursor,omitempty"`
	Resources  []ResourceTemplate `json:"resourceTemplates"`
}

type ReadResourceRequest struct {
	Uri string `json:"uri"`
}

type ReadResourceResponse struct {
	Contents []ResourceContent `json:"contents"`
}

type TextResourceContent struct {
	Uri      string  `json:"uri"`
	Text     string  `json:"text"`
	MimeType *string `json:"mimeType,omitempty"`
}

type BlobResourceContent struct {
	Uri      string  `json:"uri"`
	Blob     *string `json:"blob"`
	MimeType *string `json:"mimeType,omitempty"`
}

type ResourceContent struct {
	*TextResourceContent
	*BlobResourceContent
}

func (c ResourceContent) MarshalJSON() ([]byte, error) {
	if c.TextResourceContent != nil {
		return json.Marshal(c.TextResourceContent)
	} else if c.BlobResourceContent != nil {
		return json.Marshal(c.BlobResourceContent)
	} else {
		return []byte("{}"), nil
	}
}

// =========================================

type SetLevelRequest struct {
	Level string `json:"level"`
}

type SetLevelResponse struct {
}

// =========================================

type ServerNotification struct {
	*LoggingMessageNotification
	*ToolListChangedNotification
}

func (sn ServerNotification) ToJsonRpcNotificationInner() rpc.JsonRpcNotificationInner {
	if sn.LoggingMessageNotification != nil {
		raw, _ := json.Marshal(sn.LoggingMessageNotification)
		return rpc.JsonRpcNotificationInner{
			Method: "notifications/message",
			Params: raw,
		}
	} else if sn.ToolListChangedNotification != nil {
		raw, _ := json.Marshal(sn.ToolListChangedNotification)
		return rpc.JsonRpcNotificationInner{
			Method: "notifications/tools/list_changed",
			Params: raw,
		}
	} else {
		return rpc.JsonRpcNotificationInner{}
	}
}

type LoggingMessageNotification struct {
	Level  string `json:"level"`
	Data   string `json:"data"`
	Logger string `json:"logger,omitempty"`
}

type ToolListChangedNotification struct {
}

type McpIo interface {
	SetNotifications(notifications chan<- ServerNotification)
	Initialize(context.Context, InitializeRequest) (*InitializeResponse, error)
	ListTools(context.Context, ListToolsRequest) (*ListToolsResponse, error)
	CallTool(context.Context, CallToolRequest) (*CallToolResponse, error)
	ListPrompts(context.Context, ListPromptsRequest) (*ListPromptsResponse, error)
	GetPrompt(context.Context, GetPromptRequest) (*GetPromptResponse, error)
	ListResources(context.Context, ListResourcesRequest) (*ListResourcesResponse, error)
	ReadResource(context.Context, ReadResourceRequest) (*ReadResourceResponse, error)
	ListResourcesTemplates(context.Context, ListResourceTemplatesRequest) (*ListResourceTemplatesResponse, error)
	SetLevel(context.Context, SetLevelRequest) (*SetLevelResponse, error)
}

func wrap[x any, y any](request rpc.JsonRpcRequest, f func(context.Context, x) (*y, error)) (*rpc.JsonRpcResponse, error) {
	var ir x
	if len(request.Params) > 0 {
		if err := json.Unmarshal(request.Params, &ir); err != nil {
			slog.Error("failed to unmarshal request", slog.Any("err", err), slog.String("request", string(request.Params)))
			return nil, rpc.JsonRpcError{Code: rpc.JsonRpcInvalidParamsError, Message: err.Error(), Data: map[string]interface{}{"raw": string(request.Params)}}
		}
	}
	if irr, err := f(request.Context(), ir); err != nil {
		slog.Error("returning json rpc error", slog.Any("err", err))
		return nil, rpc.NewJsonRpcErrorFromErr(err)
	} else if raw, err := json.Marshal(irr); err != nil {
		slog.Error("failed to marshal response", slog.Any("err", err))
		return nil, rpc.NewJsonRpcErrorFromErr(err)
	} else {
		return &rpc.JsonRpcResponse{JsonRpcResponseInner: &rpc.JsonRpcResponseInner{
			Id: ref.Deref(request.Id, -1), Result: raw,
		}}, nil
	}
}

func AsHandler(inner McpIo) rpc.Handler {
	return rpc.HandlerFunc(func(req rpc.JsonRpcRequest) (*rpc.JsonRpcResponse, error) {
		switch req.Method {
		case "initialize":
			return wrap[InitializeRequest, InitializeResponse](req, inner.Initialize)
		case "tools/list":
			return wrap[ListToolsRequest, ListToolsResponse](req, inner.ListTools)
		case "tools/call":
			return wrap[CallToolRequest, CallToolResponse](req, inner.CallTool)
		case "prompts/list":
			return wrap[ListPromptsRequest, ListPromptsResponse](req, inner.ListPrompts)
		case "prompts/get":
			return wrap[GetPromptRequest, GetPromptResponse](req, inner.GetPrompt)
		case "resources/list":
			return wrap[ListResourcesRequest, ListResourcesResponse](req, inner.ListResources)
		case "resources/templates/list":
			return wrap[ListResourceTemplatesRequest, ListResourceTemplatesResponse](req, inner.ListResourcesTemplates)
		case "resources/read":
			return wrap[ReadResourceRequest, ReadResourceResponse](req, inner.ReadResource)
		case "logging/setLevel":
			return wrap[SetLevelRequest, SetLevelResponse](req, inner.SetLevel)
		default:
			if strings.HasPrefix(req.Method, "notifications/") {
				slog.Debug("dropping unsupported notification", slog.Any("method", req.Method))
				return nil, nil
			}
			slog.Warn("ignoring unknown method", slog.Any("method", req.Method))
			return nil, rpc.JsonRpcError{Code: rpc.JsonRpcMethodNotFoundError, Message: fmt.Sprintf("method not found: %s", req.Method)}
		}

	})
}
