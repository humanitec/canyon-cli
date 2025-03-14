package rpc

type Server interface {
	In() chan<- JsonRpcRequest
	Out() <-chan JsonRpcResponse
}
