package types

const (
	Op_Create_Room = 1
	Op_Join_Room   = 2
)

type Op struct {
	Type uint8
	Ts   int64
	Id   *Session
	Data interface{}
}
