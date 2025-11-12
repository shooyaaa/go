package codec

import "time"

type OpType int

const (
	Op_Create_Room OpType = 1
	Op_Join_Room   OpType = 2
	Op_Sync_Data   OpType = 3
	Op_Login       OpType = 4
	Op_Logout      OpType = 5

	//unimouse
	Op_KeyEvent   OpType = 10001
	Op_MouseEvent OpType = 10002
)

type Op struct {
	Type OpType
	Ts   int64
	Data map[string]interface{}
}

func MakeOp(op OpType, data map[string]interface{}) Op {
	return Op{Type: op, Ts: time.Now().Unix(), Data: data}
}
