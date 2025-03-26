# canyon-cli

## Installation

1. Download the latest version of the canyon-cli MCP tool from https://github.com/humanitec/canyon-cli/releases/.

2. Then execute the `install` command in the terminal: `canyon install`. This will return the JSON configuration to add to your IDE or LLM  client and will ensure the binary is executable. If security warnings pop up, you should allow or accept the binary to run.

3. Also make sure you have `humctl` installed with access to an org, hopefully the `canyon-demo` one for best results.

## Development

You can execute any of the CLI tools by running:

```
$ canyon rpc -s name=tools/list
$ canyon rpc -s name=tools/call -s arguments='{ ... }'
```

### Developing the render templates

If you're working on the HTML rendering templates, the templates are stored as the `.html.tmpl` files in the binary. 

However, it will fall back, at runtime, to the files in:

- `${HOME}/canyon-render-csv-template.html.tmpl`
- `${HOME}/canyon-render-tree-template.html.tmpl`
- `${HOME}/canyon-render-graph-template.html.tmpl`

If they are not empty. If they exist but are empty, then the default template will be written to them for development iteration.

For sample data, use the following for examples:

```
canyon rpc tools/call --stdin <<"EOF"
{
  "name": "render_csv_as_table_in_browser",
  "arguments": {
    "first_row_is_header": true,
    "raw": "id,first name,surname,age\n1,alice,berman,42\n2,bob,carren,14\n3,charles,dorito,21\n4,daphney,errol,5"
  }
}
EOF

canyon rpc tools/call --stdin <<"EOF"
{
  "name": "render_data_as_tree_in_browser",
  "arguments": {
    "root": {
      "name": "my-org",
      "data": {},
      "class": "org",
      "children": [
        {
          "name": "app-1",
          "data": {"owner": "fizz buzz"},
          "class": "app",
          "children": [
            {
              "name": "env-1",
              "data": {"name": "Environment 1"},
              "class": "env"
            },
            {
              "name": "env-2",
              "data": {"name": "Environment 2"},
              "class": "env"
            }
          ]
        }
      ]
    }
  }
}
EOF

canyon rpc tools/call --stdin <<"EOF"
{
  "name": "render_network_as_graph_in_browser",
  "arguments": {
    "nodes": [
      {
        "id": "my-org",
        "data": {},
        "class": "org",
      },
      {
        "id": "my-app",
        "data": {"owner": "fizz buzz"},
        "class": "app",
      },
      {
        "id": "env-1",
        "data": {"name": "Environment 1"},
        "class": "env",
      },
      {
        "id": "env-2",
        "data": {"name": "Environment 2"},
        "class": "env",
      }
    ],
    "links": [ 
      {"source": "my-org", "target": "my-app"},
      {"source": "my-app", "target": "env-1"},
      {"source": "my-app", "target": "env-2"}
    ]
  }
}
EOF
```
