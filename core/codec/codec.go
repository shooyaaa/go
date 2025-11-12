package codec

import (
	"encoding/json"
	"fmt"
)

type CODEC_TYPE uint8

const (
	JSON_CODEC CODEC_TYPE = iota
)

type Codec[T any] interface {
	Encode(T) ([]byte, error)
	Decode([]byte) (T, error)
}

type jsonCodec[T any] struct {
}

func (j *jsonCodec[T]) Encode(t T) ([]byte, error) {
	return json.Marshal(t)
}
func (j *jsonCodec[T]) Decode(data []byte) (T, error) {
	var t T
	err := json.Unmarshal(data, &t)
	return t, err
}

func NewCodec[T any](codec CODEC_TYPE) Codec[T] {
	switch codec {
	case JSON_CODEC:
		return &jsonCodec[T]{}
	default:
		panic(fmt.Sprintf("unknown codec type: %v", codec))
	}
}
