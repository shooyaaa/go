package actor

import (
	"context"
	"fmt"
	"strings"

	"github.com/shooyaaa/core"
)

type RpcAddress interface {
	String() string
	Transform(ctx context.Context, mail Mail[any]) *core.CoreError
}

type RpcAddressImpl struct {
	addr string
}

func (a *RpcAddressImpl) String() string {
	return a.addr
}

func (a *RpcAddressImpl) Transform(ctx context.Context, mail Mail[any]) *core.CoreError {
	return nil
}

func NewRpcAddress(addr string) RpcAddress {
	return &RpcAddressImpl{addr: addr}
}

type RpcChannelType string

const RpcChannelType_TCP RpcChannelType = "tcp"
const RpcChannelType_UDP RpcChannelType = "udp"
const RpcChannelType_IPC RpcChannelType = "ipc"
const RpcChannelType_UNIX RpcChannelType = "unix"
const RpcChannelType_FILE RpcChannelType = "file"
const RpcChannelType_PIPE RpcChannelType = "pipe"

type RpcChannel interface {
	Send(ctx context.Context, data []byte) *core.CoreError
	Receive(ctx context.Context) ([]byte, *core.CoreError)
}

type TcpChannelClient interface {
	Send(ctx context.Context, data []byte) *core.CoreError
	Receive(ctx context.Context) ([]byte, *core.CoreError)
}
type TcpChannel struct {
	addr   string
	client TcpChannelClient
}

func (c *TcpChannel) Send(ctx context.Context, data []byte) *core.CoreError {
	return nil
}

func (c *TcpChannel) Receive(ctx context.Context) ([]byte, *core.CoreError) {
	return nil, nil
}

func NewTcpChannel(addr string) RpcChannel {
	return &TcpChannel{addr: addr}
}

type HttpChannelClient interface {
	Send(ctx context.Context, data []byte) error
	Receive(ctx context.Context) ([]byte, error)
}
type HttpChannel struct {
	addr   string
	client HttpChannelClient
}

func (c *HttpChannel) Send(ctx context.Context, data []byte) *core.CoreError {
	return nil
}

func (c *HttpChannel) Receive(ctx context.Context) ([]byte, *core.CoreError) {
	if c.client != nil {
		data, err := c.client.Receive(ctx)
		if err != nil {
			return nil, core.NewCoreError(core.ERROR_CODE_MAILBOX_RECEIVE_ERROR, err.Error())
		}
		return data, nil
	}
	return nil, nil
}

func NewHttpChannel(addr string) RpcChannel {
	return &HttpChannel{addr: addr}
}

func GetChannelByAddress(addr string) RpcChannel {
	if strings.HasPrefix(addr, "http://") {
		return NewHttpChannel(addr)
	} else if strings.HasPrefix(addr, "tcp://") {
		return NewTcpChannel(addr)
	} else {
		panic(fmt.Sprintf("unimplemented address channel type: %s", addr))
	}
}
