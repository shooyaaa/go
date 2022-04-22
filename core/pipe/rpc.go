package pipe

import (
	"github.com/shooyaaa/core/network"
	"io"
)

type Rpc struct {
	network.Server
}

func (rpc *Rpc) SendTo(r io.Reader, addr Address) error {

	return nil
}

func (rpc *Rpc) ReceiveFrom(addr Address) (io.Writer, error) {

	return nil, nil
}
