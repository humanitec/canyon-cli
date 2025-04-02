package internal

import (
	"bytes"
	"encoding/json"
)

func PrettyJson(raw interface{}) string {
	buff := new(bytes.Buffer)
	enc := json.NewEncoder(buff)
	enc.SetIndent("", "  ")
	_ = enc.Encode(raw)
	return buff.String()
}
