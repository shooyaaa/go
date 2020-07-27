package types

import (
	"encoding/json"
)

type Json struct {
}

func (j *Json) Encode(i []Op) ([]byte, error) {
	return json.Marshal(i)
}

func (j *Json) Decode(data []byte, i *[]Op) (int, error) {
	err := json.Unmarshal(data, i)
	return len(data), err
}
