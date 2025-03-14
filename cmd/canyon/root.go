package main

import (
	"fmt"
	"io"
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

		var w io.Writer = cmd.ErrOrStderr()
		if lf, _ := cmd.Flags().GetString("log-file"); lf != "" {
			logf, err := os.OpenFile(lf, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, fmt.Errorf("failed to open log file: %v", err))
				os.Exit(1)
			}
			w = io.MultiWriter(logf, cmd.ErrOrStderr())
			go func() {
				<-cmd.Context().Done()
				_ = logf.Close()
			}()
		}
		internal.SetupLogging(d, w)
		return nil
	},
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s %s", internal.ModulePath, internal.ModuleVersion)
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Increase log verbosity to debug level")
	rootCmd.PersistentFlags().String("log-file", "", "Direct structured logging output to the given log file rather than stderr")
}
