package rpc

import (
	"encoding/json"
	"fmt"
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

func (c ServerNotification) MarshalJSON() ([]byte, error) {
	if c.LoggingMessageNotification != nil {
		return json.Marshal(c.LoggingMessageNotification)
	} else if c.ToolListChangedNotification != nil {
		return json.Marshal(c.ToolListChangedNotification)
	} else {
		return []byte("{}"), nil
	}
}

type LoggingMessageNotificationMethod struct {
}

func (t LoggingMessageNotificationMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal("notification/message")
}

type LoggingMessageNotification struct {
	Method LoggingMessageNotificationMethod `json:"method"`
	Params LoggingMessageParams             `json:"params"`
}

type LoggingMessageParams struct {
	Level  string `json:"level"`
	Data   string `json:"data"`
	Logger string `json:"logger,omitempty"`
}

type ToolListChangedNotificationMethod struct {
}

func (t ToolListChangedNotificationMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal("notification/tools/list_changed")
}

type ToolListChangedNotification struct {
	Method ToolListChangedNotificationMethod `json:"method"`
	Params map[string]interface{}            `json:"params"`
}
