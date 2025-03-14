package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
)

type JsonRpcErrorCode int

const (
	JsonRpcNotFound            JsonRpcErrorCode = -32002
	JsonRpcParseError          JsonRpcErrorCode = -32700
	JsonRpcInvalidRequestError JsonRpcErrorCode = -32600
	JsonRpcMethodNotFoundError JsonRpcErrorCode = -32601
	JsonRpcInvalidParamsError  JsonRpcErrorCode = -32602
	JsonRpcInternalError       JsonRpcErrorCode = -32603
)

type JsonRpcRequest struct {
	ctx     context.Context
	JsonRpc JsonRpcVersion  `json:"jsonrpc"`
	Id      *int            `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func (j JsonRpcRequest) Context() context.Context {
	if j.ctx == nil {
		return context.TODO()
	}
	return j.ctx
}

func (j JsonRpcRequest) WithContext(ctx context.Context) JsonRpcRequest {
	j.ctx = ctx
	return j
}

func (j JsonRpcRequest) LogValue() slog.Value {
	if raw, err := json.Marshal(j); err != nil {
		return slog.StringValue(fmt.Sprintf("%+v", err))
	} else {
		return slog.StringValue(string(raw))
	}
}

var _ slog.LogValuer = JsonRpcRequest{}

type JsonRpcResponse struct {
	ctx     context.Context
	JsonRpc JsonRpcVersion `json:"jsonrpc"`
	*JsonRpcResponseInner
	*JsonRpcNotificationInner
}

type JsonRpcResponseInner struct {
	Id     int             `json:"id,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *JsonRpcError   `json:"error,omitempty"`
}

type JsonRpcNotification interface {
	ToJsonRpcNotificationInner() JsonRpcNotificationInner
}

type JsonRpcNotificationInner struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

func (j JsonRpcResponse) Context() context.Context {
	if j.ctx == nil {
		return context.TODO()
	}
	return j.ctx
}

func (j JsonRpcResponse) WithContext(ctx context.Context) JsonRpcResponse {
	j.ctx = ctx
	return j
}

func (j JsonRpcResponse) LogValue() slog.Value {
	if raw, err := json.Marshal(j); err != nil {
		return slog.StringValue(fmt.Sprintf("%+v", err))
	} else {
		return slog.StringValue(string(raw))
	}
}

var _ slog.LogValuer = JsonRpcResponse{}

type JsonRpcError struct {
	Code    JsonRpcErrorCode       `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (err JsonRpcError) Error() string {
	return fmt.Sprintf("json rpc error: %d: %s (%#v)", err.Code, err.Message, err.Data)
}

func NewJsonRpcErrorFromErr(err error) JsonRpcError {
	if e := (*JsonRpcError)(nil); errors.As(err, &e) {
		return *e
	}
	return JsonRpcError{Code: JsonRpcInternalError, Message: err.Error()}
}
