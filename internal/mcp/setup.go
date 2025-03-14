package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func New() McpIo {
	m := Impl{
		Instructions: "",
		Tools:        make([]Tool, 0),
	}

	t := Tool{
		Name: "list-canyon-paths",
		Description: `Returns a list of 'paths' supported by the canyon MCP server.
Paths are remote functions which can be used to query or achieve a wide array of functionality.
The list of available paths may change over time so consider listing the available paths when there is low confidence that an existing paths can be used to solve the user query`,
		InputSchema: map[string]interface{}{"type": "object"},
		Callable: func(ctx context.Context, arguments map[string]interface{}) ([]CallToolResponseContent, error) {
			raw, err := json.Marshal([]ToolResponse{
				{Name: "get-bananaboats", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"count": map[string]interface{}{"type": "integer"}}, "required": []interface{}{"count"}}},
				{Name: "get-skyhooks"},
				{Name: "wait-seconds", InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{"count": map[string]interface{}{"type": "integer"}}, "required": []interface{}{"count"}}},
			})
			if err != nil {
				return nil, err
			}
			return []CallToolResponseContent{
				NewTextToolResponseContent("Here's an array of the current canyon tools in JSON: " + string(raw)),
			}, nil
		},
	}
	t2 := Tool{
		Name:        "call-canyon-path",
		Description: "Call a canyon path previously discovered through list-canyon-paths",
		InputSchema: map[string]interface{}{"type": "object", "properties": map[string]interface{}{
			"name":      map[string]interface{}{"type": "string"},
			"arguments": map[string]interface{}{"type": "object"},
		}, "required": []interface{}{"name"}},
		Callable: func(ctx context.Context, arguments map[string]interface{}) ([]CallToolResponseContent, error) {
			name, _ := arguments["name"].(string)
			args, _ := arguments["arguments"].(map[string]interface{})
			switch name {
			case "get-bananaboats":
				return []CallToolResponseContent{
					NewTextToolResponseContent("The path is being executed. please use the get-skyhooks command to monitor it's progress."),
				}, nil
			case "get-skyhooks":
				return []CallToolResponseContent{
					NewTextToolResponseContent("The path is still being executed. please query again later to check if it has completed"),
				}, nil
			case "wait-seconds":
				var seconds int64
				switch typed := args["count"].(type) {
				case float64:
					seconds = int64(typed)
				case int64:
					seconds = typed
				case int:
					seconds = int64(typed)
				case int32:
					seconds = int64(typed)
				}
				time.Sleep(time.Duration(seconds) * time.Second)
				return []CallToolResponseContent{NewTextToolResponseContent(fmt.Sprintf("done %s - please use the 'check-wait' path to confirm the duration", time.Duration(seconds)*time.Second))}, nil
			default:
				return []CallToolResponseContent{NewTextToolResponseContent("unknown path name")}, nil
			}
		},
	}
	m.InjectTools(t, t2)
	return &m
}
