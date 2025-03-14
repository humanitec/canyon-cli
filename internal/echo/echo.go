package echo

import (
	"log/slog"

	"github.com/humanitec/canyon-cli/internal/mcp"
	"github.com/humanitec/canyon-cli/internal/ref"
	"github.com/humanitec/canyon-cli/internal/rpc"
)

func NewEchoServer() rpc.Server {
	return &rpc.Generic{
		Handler: rpc.HandlerFunc(func(req rpc.JsonRpcRequest) (*rpc.JsonRpcResponse, error) {
			slog.Debug("Echoing request", slog.Any("req", req.LogValue()))

			if c := rpc.GetNotificationChannel(req.Context()); c != nil {
				c <- mcp.ServerNotification{LoggingMessageNotification: &mcp.LoggingMessageNotification{
					Level: "info",
					Data:  "this is a log message",
				}}
			}

			return ref.Ref(rpc.JsonRpcResponse{
				JsonRpcResponseInner: &rpc.JsonRpcResponseInner{
					Id: ref.Deref(req.Id, -1),
					Error: &rpc.JsonRpcError{
						Code:    rpc.JsonRpcMethodNotFoundError,
						Message: "method not found",
						Data: map[string]interface{}{
							"method": req.Method,
							"params": req.Params,
						},
					},
				},
			}.WithContext(req.Context())), nil
		}),
	}
}
