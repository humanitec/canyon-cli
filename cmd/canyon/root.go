package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/humanitec/canyon-cli/internal"
)

var rootCmd = &cobra.Command{
	Use:          "canyon",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s %s", internal.ModulePath, internal.ModuleVersion)
}
