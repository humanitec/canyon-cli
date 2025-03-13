package rpc

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToolContentEncoding(t *testing.T) {
	v := TextContentType{}
	raw, _ := json.Marshal(v)
	assert.Equal(t, "\"text\"", string(raw))

	x := TextContent{Text: "something"}
	raw, _ = json.Marshal(x)
	assert.Equal(t, "{\"type\":\"text\",\"text\":\"something\"}", string(raw))

	o := NewTextToolResponseContentWithAudience("something", "aud")
	raw, _ = json.Marshal(o)
	assert.Equal(t, "{\"type\":\"text\",\"text\":\"something\",\"annotations\":{\"audience\":[\"aud\"]}}", string(raw))
}
