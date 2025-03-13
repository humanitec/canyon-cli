package rpc

import "encoding/json"

type JsonRpcVersion struct {
}

func (r JsonRpcVersion) MarshalJSON() ([]byte, error) {
	return json.Marshal("2.0")
}

func (r JsonRpcVersion) UnmarshalJSON(b []byte) error {
	return nil
}
