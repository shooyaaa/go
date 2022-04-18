package types

type Codec interface {
	Encode([]Op) ([]byte, error)
	Decode([]byte, *[]Op) (int, error)
}
