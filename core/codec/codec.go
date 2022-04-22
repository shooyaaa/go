package codec

type CODEC_TYPE uint8

const (
	JSON_CODEC CODEC_TYPE = iota
)

type Codec interface {
	Encode(interface{}) ([]byte, error)
	Decode([]byte) (interface{}, int, error)
}

func GetCodec(codec CODEC_TYPE) Codec {
	switch codec {
	case JSON_CODEC:
		return &Json{}
	}
	return nil
}
