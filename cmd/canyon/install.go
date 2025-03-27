package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:           "install",
	Short:         "Show the installation configuration for clients",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		h, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cannot find the user home directory, you'll need to manage the installation manually")
		}

		command := os.Args[0]
		if p, _ := filepath.Abs(command); p != "" {
			command = p
		}

		mcpServerConfig := map[string]interface{}{
			"mcpServers": map[string]interface{}{
				"canyon": map[string]interface{}{
					"command": command,
					"args": []interface{}{
						"mcp",
					},
					"env": map[string]interface{}{
						"HOME": h,
					},
				},
			},
		}
		buff := new(bytes.Buffer)
		enc := json.NewEncoder(buff)
		enc.SetIndent("", "  ")
		_ = enc.Encode(mcpServerConfig)
		raw := buff.String()

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), `Configure the following MCP server in your LLM Client:

%s

If you are using 'Cline' VSCode extension, configure this in Cline extension > MCP Servers > cline_mcp_settings.json.

If you are using 'Claude Desktop', configure this in the ~/Library/Application Support/Claude/claude_desktop_config.json file or by going to the developer settings pane.

If you are using another Client, google for 'configure mcp servers in XYZ'.
`, raw)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
