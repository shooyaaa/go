package session

import (
	"bytes"
	"encoding/json"
)

type Json struct {
}

func (j *Json) Encode(o Op) ([]byte, error) {
	return json.Marshal(o)
}

func (j *Json) Decode(data []byte) (*Op, int, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	var o *Op = &Op{}
	for dec.More() {
		err := dec.Decode(o)
		if err != nil {
			return nil, 0, err
		} else {
			break
		}
	}
	return o, int(dec.InputOffset()), nil
}
