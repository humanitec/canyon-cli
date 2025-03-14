package rpc

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime/debug"

	"github.com/humanitec/canyon-cli/internal/ref"
)

func LoggingMiddleware(next Handler) Handler {
	return HandlerFunc(func(req JsonRpcRequest) (*JsonRpcResponse, error) {
		logger := slog.Default().With(slog.Int("id", ref.Deref(req.Id, -1)))
		logger.Debug("received", slog.Any("req", req.LogValue()))
		res, err := next.Handle(req)
		if err != nil {
			var e JsonRpcError
			if !errors.As(err, &e) || e.Code == JsonRpcInternalError {
				logger.Error("sending", slog.Any("err", err.Error()))
			} else {
				logger.Debug("sending", slog.Any("err", err.Error()))
			}
		} else if res != nil {
			logger.Debug("sending", slog.Any("res", res.LogValue()))
		}
		return res, err
	})
}

var _ MiddlewareFunc = LoggingMiddleware

func RecoveryMiddleware(next Handler) Handler {
	return HandlerFunc(func(req JsonRpcRequest) (res *JsonRpcResponse, err error) {
		defer func() {
			if e := recover(); e != nil {
				res = nil
				err = fmt.Errorf("caught panic: %v at:\n%s", e, string(debug.Stack()))
			}
		}()
		res, err = next.Handle(req)
		return
	})
}
