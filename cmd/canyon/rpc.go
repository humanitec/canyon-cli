package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/humanitec/canyon-cli/internal/mcp"
	"github.com/humanitec/canyon-cli/internal/mcp/tools"
	"github.com/humanitec/canyon-cli/internal/ref"
	"github.com/humanitec/canyon-cli/internal/rpc"
)

var rpcCmd = &cobra.Command{
	Use:           "rpc",
	Args:          cobra.ExactArgs(1),
	Short:         "Send an individual RPC message and observe the results.",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		intermediate := make(map[string]interface{})
		if b, _ := cmd.Flags().GetBool("stdin"); b {
			dec := yaml.NewDecoder(cmd.InOrStdin())
			if err := dec.Decode(&intermediate); err != nil {
				return fmt.Errorf("failed to decode JSON map from stdin: %w", err)
			}
		}

		rawParams, _ := cmd.Flags().GetStringToString("set")
		for k, rv := range rawParams {
			if rv == "" {
				delete(intermediate, k)
			} else {
				var v interface{}
				if err := json.Unmarshal([]byte(rv), &v); err != nil {
					intermediate[k] = rv
				} else {
					intermediate[k] = v
				}
			}
		}
		requestId := int(rand.Int64())
		rawRawParams, _ := json.Marshal(intermediate)
		slog.Info("executing method with params", slog.String("method", args[0]), slog.String("params", string(rawRawParams)), slog.Int("request_id", requestId))

		h := mcp.AsHandler(tools.New())
		h = rpc.RecoveryMiddleware(h)
		h = rpc.LoggingMiddleware(h)
		server := &rpc.Generic{Handler: h}
		in := server.In()
		defer close(in)

		go func() {
			in <- rpc.JsonRpcRequest{
				Method: args[0],
				Id:     ref.Ref(requestId),
				Params: rawRawParams,
			}.WithContext(cmd.Context())
		}()

		out := server.Out()
		for {
			select {
			case result := <-out:
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				if err := enc.Encode(result); err != nil {
					return err
				}
				if result.JsonRpcResponseInner != nil && result.JsonRpcResponseInner.Id == requestId {
					return nil
				}
			case <-cmd.Context().Done():
				return cmd.Context().Err()
			}
		}
	},
}

func init() {
	rpcCmd.Flags().StringToStringP("set", "s", nil, "Set key-value params")
	rpcCmd.Flags().Bool("stdin", false, "Read params from stdin")
	rootCmd.AddCommand(rpcCmd)
}
