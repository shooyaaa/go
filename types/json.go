package types

import (
	"encoding/json"
)

type Json struct {
}

func (j *Json) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j *Json) Decode(data []byte, i interface{}) (int, error) {
	err := json.Unmarshal(data, i)
	return len(data), err
}
