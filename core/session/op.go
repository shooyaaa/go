package session

import "time"

type OpType int

const (
	Op_Create_Room OpType = 1
	Op_Join_Room          = 2
	Op_Sync_Data          = 3
	Op_Login              = 4
	Op_Logout             = 4

	//unimouse
	Op_KeyEvent   = 10001
	Op_MouseEvent = 10002
)

type Op struct {
	Type OpType
	Ts   int64
	Data map[string]interface{}
}

func MakeOp(op OpType, data map[string]interface{}) Op {
	return Op{Type: op, Ts: time.Now().Unix(), Data: data}
}
