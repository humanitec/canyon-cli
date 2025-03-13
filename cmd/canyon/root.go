package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/humanitec/canyon-cli/internal"
)

var rootCmd = &cobra.Command{
	Use:          "canyon",
	SilenceUsage: true,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		d, _ := cmd.Flags().GetBool("debug")
		d = d || strings.ToLower(os.Getenv("CANYON_CLI_DEBUG")) == "true"
		internal.SetupLogging(d, cmd.ErrOrStderr())
		return nil
	},
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s %s", internal.ModulePath, internal.ModuleVersion)
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Increase log verbosity to debug level")
}
