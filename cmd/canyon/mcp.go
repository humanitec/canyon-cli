package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/spf13/cobra"

	"github.com/humanitec/canyon-cli/internal/mcp"
	"github.com/humanitec/canyon-cli/internal/rpc"
)

var mcpCmd = &cobra.Command{
	Use:           "mcp",
	Short:         "Start the raw stdio mcp session normally used by LLM clients.",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		h := mcp.AsHandler(mcp.New())
		h = rpc.RecoveryMiddleware(h)
		h = rpc.LoggingMiddleware(h)
		server := &rpc.Generic{Handler: h}
		in := server.In()

		scanner := bufio.NewScanner(cmd.InOrStdin())
		errChan := make(chan error)
		go func() {
			defer func() {
				slog.Info("Closing input session")
				close(in)
			}()
			for {
				select {
				case <-cmd.Context().Done():
					return
				case <-time.After(time.Millisecond * 100):
				}
				for scanner.Scan() {
					if len(scanner.Bytes()) == 0 {
						break
					}
					var msg rpc.JsonRpcRequest
					dec := json.NewDecoder(bytes.NewReader(scanner.Bytes()))
					dec.DisallowUnknownFields()
					if err := dec.Decode(&msg); err != nil {
						errChan <- fmt.Errorf("failed to read json formatted line '%q' as a request: %w", scanner.Text(), err)
						return
					}
					server.In() <- msg.WithContext(cmd.Context())
				}
			}
		}()

		for {
			select {
			case err := <-errChan:
				return err
			case <-cmd.Context().Done():
				return cmd.Context().Err()
			case r := <-server.Out():
				enc := json.NewEncoder(cmd.OutOrStdout())
				if err := enc.Encode(r); err != nil {
					return fmt.Errorf("failed to encode response: %w", err)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
