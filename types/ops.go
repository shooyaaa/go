package types

type OpCode uint16

type OpPing struct {
	Op   string
	Data int64
}

type Dispatcher interface {
	Handle(OpCode, []byte) error
}
