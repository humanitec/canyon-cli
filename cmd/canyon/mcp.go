package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/humanitec/canyon-cli/internal/rpc"
)

var mcpCmd = &cobra.Command{
	Use:           "mcp",
	Short:         "Start the raw stdio mcp session normally used by LLM clients.",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		server := rpc.NewEchoServer()
		in := server.In()
		defer close(in)

		scanner := bufio.NewScanner(cmd.InOrStdin())
		errChan := make(chan error)
		go func() {
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
					if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
						errChan <- fmt.Errorf("failed to read json formatted line '%q' as a request: %w", scanner.Text(), err)
						return
					}
					select {
					case server.In() <- msg:
					default:
						return
					}
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
