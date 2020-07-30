package types

type Conn interface {
	Read() ([]byte, error)
	Write(bytes []byte) (int, error)
}
