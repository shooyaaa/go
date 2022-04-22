package codec

import (
	"bytes"
	"encoding/json"
)

type Json struct {
}

func (j *Json) Encode(i interface{}) ([]byte, error) {
	return json.Marshal(i)
}

func (j *Json) Decode(data []byte) (interface{}, int, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	var i interface{}
	for dec.More() {
		err := dec.Decode(i)
		if err != nil {
			return nil, 0, err
		} else {
			break
		}
	}
	return i, int(dec.InputOffset()), nil
}
