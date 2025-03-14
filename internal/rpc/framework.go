package rpc

import (
	"context"
	"sync"

	"github.com/humanitec/canyon-cli/internal/ref"
)

type Server interface {
	In() chan<- JsonRpcRequest
	Out() <-chan JsonRpcResponse
}

type Handler interface {
	Handle(req JsonRpcRequest, notifications chan<- JsonRpcNotification) (*JsonRpcResponse, error)
}

type HandlerFunc func(req JsonRpcRequest, notifications chan<- JsonRpcNotification) (*JsonRpcResponse, error)

func (f HandlerFunc) Handle(req JsonRpcRequest, notifications chan<- JsonRpcNotification) (*JsonRpcResponse, error) {
	return f(req, notifications)
}

type Generic struct {
	Handler Handler

	in   chan JsonRpcRequest
	out  chan JsonRpcResponse
	once sync.Once
}

func (e *Generic) In() chan<- JsonRpcRequest {
	e.setup()
	return e.in
}

func (e *Generic) Out() <-chan JsonRpcResponse {
	e.setup()
	return e.out
}

func (e *Generic) setup() {
	e.once.Do(func() {
		e.in = make(chan JsonRpcRequest)
		e.out = make(chan JsonRpcResponse)

		notifications := make(chan JsonRpcNotification)
		notificationCtx, notificationsCancel := context.WithCancel(context.Background())
		go func() {
			for {
				select {
				case <-notificationCtx.Done():
					return
				case n := <-notifications:
					e.out <- JsonRpcResponse{
						JsonRpcNotificationInner: ref.Ref(n.ToJsonRpcNotificationInner()),
					}
				}
			}
		}()

		go func() {
			defer notificationsCancel()
			for req := range e.in {
				r, err := e.Handler.Handle(req, notifications)
				if err != nil {
					r = ref.Ref(JsonRpcResponse{
						JsonRpcResponseInner: &JsonRpcResponseInner{
							Id: ref.Deref(req.Id, -1),
							Error: &JsonRpcError{
								Code:    JsonRpcInternalError,
								Message: "internal error",
								Data: map[string]interface{}{
									"message": err.Error(),
								},
							},
						},
					}.WithContext(r.Context()))
				}
				e.out <- *r
			}
		}()
	})
}
