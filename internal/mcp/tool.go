package mcp

import "context"

type Tool struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
	Callable    func(ctx context.Context, arguments map[string]interface{}) ([]CallToolResponseContent, error)
}
