package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

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
		rawRawParams, _ := json.Marshal(intermediate)

		server := rpc.NewEchoServer()
		in := server.In()
		defer close(in)

		requestId := int(rand.Int64())
		go func() {
			in <- rpc.JsonRpcRequest{
				Method: args[0],
				Id:     requestId,
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
				if result.Id != nil && *result.Id == requestId {
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
