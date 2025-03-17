package tools

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"strings"

	"github.com/pkg/browser"

	"github.com/humanitec/canyon-cli/internal/mcp"
)

// RenderCSVAsTable renders csv as a table.
func NewRenderCSVAsTable() mcp.Tool {
	tmpl, err := template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
</head>
<body>
	<table>
	{{ with .header }}
	<thead>
		<tr>
			{{ range .}}
			<th scope="col">{{.}}</th>	
			{{ end}}
		</tr>
	</thead>
	{{ end }}
	<tbody>
		{{ range .rows}}
		<tr>
			{{ range . }}
			<td>{{ . }}</td>
			{{ end }}
		</tr>
		{{ end }}
	</tbody>
	</table>
</body>
</html>
`)
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

// RenderTreeAsTree renders a hierarchy as basic html. Use https://d3js.org/d3-hierarchy/hierarchy instead in the future.
func NewRenderTreeAsTree() mcp.Tool {
	tmpl, err := template.New("").Parse(`
{{define "T"}}
	{{ .name }}
	<ol>
	{{ range .children}}
	<li>{{ template "T" . }}</li>
	{{ end }}
	</ol>
{{end}}

<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
</head>
<body>
	{{ template "T" .root }}
</body>
</html>
`)
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
						"name":     map[string]interface{}{"type": "string"},
						"children": map[string]interface{}{"type": "array", "items": map[string]interface{}{"$ref": "#/$defs/node"}},
					},
					"required": []interface{}{"name"},
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
	tmpl, err := template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
</head>
<body>
	<pre><code>{{ .data }}</code></pre>
</body>
</html>
`)
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
						"group": map[string]interface{}{"type": "integer"},
					},
					"required": []interface{}{"id"},
				}},
				"links": map[string]interface{}{"type": "array", "description": "The list of links between nodes in the network", "items": map[string]interface{}{
					"type":        "object",
					"description": "A link in the network graph",
					"properties": map[string]interface{}{
						"source": map[string]interface{}{"type": "string"},
						"target": map[string]interface{}{"type": "string"},
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
