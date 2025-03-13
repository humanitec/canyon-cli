package rpc

import (
	"log/slog"
	"sync"

	"github.com/humanitec/canyon-cli/internal/ref"
)

type Server interface {
	In() chan<- JsonRpcRequest
	Out() <-chan JsonRpcResponse
}

type echoServer struct {
	in   chan JsonRpcRequest
	out  chan JsonRpcResponse
	once sync.Once
}

func (e *echoServer) In() chan<- JsonRpcRequest {
	e.setup()
	return e.in
}

func (e *echoServer) Out() <-chan JsonRpcResponse {
	e.setup()
	return e.out
}

func (e *echoServer) setup() {
	e.once.Do(func() {
		e.in = make(chan JsonRpcRequest)
		e.out = make(chan JsonRpcResponse)
		go func() {
			for req := range e.in {
				slog.Debug("Echoing request", slog.Any("req", req))
				e.out <- JsonRpcResponse{
					Id: ref.Ref(req.Id),
					Error: &JsonRpcError{
						Code:    JsonRpcMethodNotFoundError,
						Message: "method not found",
						Data: map[string]interface{}{
							"method": req.Method,
							"params": req.Params,
						},
					},
				}
			}
		}()
	})
}

func NewEchoServer() Server {
	return new(echoServer)
}
