package session

import (
	"bytes"
	"encoding/json"

	"github.com/shooyaaa/core/codec"
)

type Json struct {
}

func (j *Json) Encode(o codec.Op) ([]byte, error) {
	return json.Marshal(o)
}

func (j *Json) Decode(data []byte) (codec.Op, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	var o codec.Op
	err := dec.Decode(&o)
	return o, err
}
