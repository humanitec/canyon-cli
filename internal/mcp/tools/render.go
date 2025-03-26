package tools

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/browser"

	"github.com/humanitec/canyon-cli/internal/mcp"
)

//go:embed render_csv.html.tmpl
var renderCsvTemplate string

//go:embed render_tree.html.tmpl
var renderTreeTemplate string

//go:embed render_graph.html.tmpl
var renderGraphTemplate string

func init() {

	f := func(path string, defaultContent string) string {
		raw, err := os.ReadFile(path)
		if err == nil {
			if len(raw) == 0 {
				_ = os.WriteFile(path, []byte(defaultContent), 0644)
			} else {
				return string(raw)
			}
		}
		return defaultContent
	}

	h, err := os.UserHomeDir()
	if err == nil {
		renderCsvTemplate = f(filepath.Join(h, "canyon-render-csv-template.html.tmpl"), renderCsvTemplate)
		renderTreeTemplate = f(filepath.Join(h, "canyon-render-tree-template.html.tmpl"), renderTreeTemplate)
		renderGraphTemplate = f(filepath.Join(h, "canyon-render-graph-template.html.tmpl"), renderGraphTemplate)
	}
}

// NewRenderCSVAsTable renders csv as a table.
func NewRenderCSVAsTable() mcp.Tool {
	tmpl, err := template.New("").Parse(renderCsvTemplate)
	if err != nil {
		panic(err)
	}
	return mcp.Tool{
		Name:        "render_csv_as_table_in_browser",
		Description: `This tool can be used to render CSV data as a pretty html table in the users browser`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"raw":                 map[string]interface{}{"type": "string", "description": "The raw multiline csv content"},
				"first_row_is_header": map[string]interface{}{"type": "boolean", "description": "Whether the first row of csv is the header"},
			},
			"required": []interface{}{"raw"},
		},
		Callable: func(ctx context.Context, m map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			r := csv.NewReader(strings.NewReader(m["raw"].(string)))
			rows, err := r.ReadAll()
			if err != nil {
				return nil, fmt.Errorf("invalid csv content")
			}
			hasHeader, _ := m["first_row_is_header"].(bool)
			var header []string
			if hasHeader {
				header = rows[0]
				rows = rows[1:]
			}
			buffer := new(bytes.Buffer)
			if err := tmpl.Execute(buffer, map[string]interface{}{
				"header": header,
				"rows":   rows,
			}); err != nil {
				slog.Error("failed to execute template", slog.Any("err", err))
				return nil, fmt.Errorf("could not render html content")
			}
			if err := browser.OpenReader(bytes.NewReader(buffer.Bytes())); err != nil {
				return nil, err
			}
			return []mcp.CallToolResponseContent{mcp.NewTextToolResponseContent("browser was opened")}, nil
		},
	}
}

// NewRenderTreeAsTree renders a hierarchy as basic html. Use https://d3js.org/d3-hierarchy/hierarchy instead in the future.
func NewRenderTreeAsTree() mcp.Tool {
	tmpl, err := template.New("").Parse(renderTreeTemplate)
	if err != nil {
		panic(err)
	}

	return mcp.Tool{
		Name:        "render_data_as_tree_in_browser",
		Description: `This tool will render a hierarchy such as a tree structure or directory tree in a user friendly way in the browser.`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"root": map[string]interface{}{"$ref": "#/$defs/node", "description": "The root of the tree structure"},
			},
			"required": []interface{}{"root"},
			"$defs": map[string]interface{}{
				"node": map[string]interface{}{
					"type":        "object",
					"description": "A node in the tree structure",
					"properties": map[string]interface{}{
						"name":     map[string]interface{}{"type": "string", "description": "The name of the node"},
						"class":    map[string]interface{}{"type": "string", "description": "The class of the node. Well known classes are: 'org', 'app', 'env_type', 'env', 'workload', 'resource', and 'other' but arbitrary strings can be used too"},
						"data":     map[string]interface{}{"type": "object", "description": "Arbitrary additional metadata to include on the node visualisation"},
						"children": map[string]interface{}{"type": "array", "items": map[string]interface{}{"$ref": "#/$defs/node"}},
					},
					"required": []interface{}{"name", "class"},
				},
			},
		},
		Callable: func(ctx context.Context, arguments map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			root, _ := arguments["root"].(map[string]interface{})
			buffer := new(bytes.Buffer)
			if err := tmpl.Execute(buffer, map[string]interface{}{
				"root": root,
			}); err != nil {
				slog.Error("failed to execute template", slog.Any("err", err))
				return nil, fmt.Errorf("could not render html content")
			}
			if err := browser.OpenReader(bytes.NewReader(buffer.Bytes())); err != nil {
				return nil, err
			}
			return []mcp.CallToolResponseContent{mcp.NewTextToolResponseContent("browser was opened")}, nil
		},
	}
}

func NewRenderNetworkAsGraph() mcp.Tool {
	tmpl, err := template.New("").Parse(renderGraphTemplate)
	if err != nil {
		panic(err)
	}
	return mcp.Tool{
		Name:        "render_network_as_graph_in_browser",
		Description: `This tool will render an interconnected network as a force directed graph in the browser.`,
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"nodes": map[string]interface{}{"type": "array", "description": "The list of nodes in the network", "items": map[string]interface{}{
					"type":        "object",
					"description": "A node in the network graph",
					"properties": map[string]interface{}{
						"id":    map[string]interface{}{"type": "string"},
						"class": map[string]interface{}{"type": "string", "description": "The class of the node. Well known classes are: 'org', 'app', 'env_type', 'env', 'workload', 'resource', and 'other' but arbitrary strings can be used too"},
						"data":  map[string]interface{}{"type": "object", "description": "Arbitrary additional metadata to include on the node visualisation"},
					},
					"required": []interface{}{"id", "class"},
				}},
				"links": map[string]interface{}{"type": "array", "description": "The list of links between nodes in the network", "items": map[string]interface{}{
					"type":        "object",
					"description": "A link in the network graph",
					"properties": map[string]interface{}{
						"source": map[string]interface{}{"type": "string", "description": "The source node id of the link"},
						"target": map[string]interface{}{"type": "string", "description": "The target node id of the link"},
					},
					"required": []interface{}{"source", "target"},
				}},
			},
			"required": []interface{}{"nodes", "links"},
		},
		Callable: func(ctx context.Context, arguments map[string]interface{}) ([]mcp.CallToolResponseContent, error) {
			raw, _ := json.MarshalIndent(arguments, "", "  ")
			buffer := new(bytes.Buffer)
			if err := tmpl.Execute(buffer, map[string]interface{}{
				"data": string(raw),
			}); err != nil {
				slog.Error("failed to execute template", slog.Any("err", err))
				return nil, fmt.Errorf("could not render html content")
			}
			if err := browser.OpenReader(bytes.NewReader(buffer.Bytes())); err != nil {
				return nil, err
			}
			return []mcp.CallToolResponseContent{mcp.NewTextToolResponseContent("browser was opened")}, nil
		},
	}
}
