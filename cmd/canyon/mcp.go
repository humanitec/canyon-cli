package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:           "mcp",
	Short:         "Start the raw stdio mcp session normally used by LLM clients.",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return fmt.Errorf("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
