package types

const (
	Op_Create_Room = 1
	Op_Join_Room   = 2
	Op_Sync_Data   = 3
	Op_Login       = 4
	Op_Logout      = 4
)

type Op struct {
	Type uint8
	Ts   int64
	Data map[string]float64
	id   *Session
}

func (o *Op) SetId(id *Session) {
	o.id = id
}

func (o *Op) GetId() *Session {
	return o.id
}
