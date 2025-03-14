package echo

import (
	"log/slog"

	"github.com/humanitec/canyon-cli/internal/mpc"
	"github.com/humanitec/canyon-cli/internal/ref"
	"github.com/humanitec/canyon-cli/internal/rpc"
)

func NewEchoServer() rpc.Server {
	return &rpc.Generic{
		Handler: rpc.HandlerFunc(func(req rpc.JsonRpcRequest, notifications chan<- rpc.JsonRpcNotification) (*rpc.JsonRpcResponse, error) {
			slog.Debug("Echoing request", slog.Any("req", req.LogValue()))

			notifications <- mpc.ServerNotification{LoggingMessageNotification: &mpc.LoggingMessageNotification{
				Level: "info",
				Data:  "this is a log message",
			}}

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
