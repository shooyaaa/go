package pipe

import "io"

type Address interface {
	Desc() interface{}
}

type Pipe interface {
	SendTo(io.Reader, Address) error
	ReceiveFrom(Address) (io.Writer, error)
}
