package types

type Op struct {
	Type uint8
	Ts 	 int64
	Data interface{}
}

type Dispatcher interface {
}
