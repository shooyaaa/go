package types

import (
	"bytes"
	"encoding/json"
)

type Json struct {
}

func (j *Json) Encode(i []Op) ([]byte, error) {
	return json.Marshal(i)
}

func (j *Json) Decode(data []byte, i *[]Op) (int, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	for dec.More() {
		err := dec.Decode(i)
		if err != nil {
			return 0, err
		} else {
			break
		}
	}
	return int(dec.InputOffset()), nil
}
