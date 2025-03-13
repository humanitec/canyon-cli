package rpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSONRpcVersion(t *testing.T) {
	o := JsonRpcVersion{}
	raw, _ := json.Marshal(o)
	assert.Equal(t, "\"2.0\"", string(raw))

	x := JsonRpcResponse{}
	raw, _ = json.Marshal(x)
	assert.Equal(t, "{\"jsonrpc\":\"2.0\"}", string(raw))
}
