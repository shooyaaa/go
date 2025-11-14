package actor

import (
	"context"
	"fmt"

	"github.com/shooyaaa/core"
	"github.com/shooyaaa/core/codec"
)

type AddressType string

const AddressType_LOCAL AddressType = "local"
const AddressType_REMOTE AddressType = "remote"

type Address interface {
	String() string
	Transform(ctx context.Context, mail Mail[any]) *core.CoreError
}

type LocalAddress struct {
	postman Postman
}

func (a *LocalAddress) String() string {
	return fmt.Sprintf("%s:%s", AddressType_LOCAL, a.postman.ID())
}

func (a *LocalAddress) Transform(ctx context.Context, mail Mail[any]) *core.CoreError {
	return a.postman.Receive(ctx, mail)
}
func NewLocalPostManAddress(postman Postman) Address {
	return &LocalAddress{postman: postman}
}

type RemoteAddress struct {
	address RpcAddress
}

func (a *RemoteAddress) String() string {
	return fmt.Sprintf("%s:%s", AddressType_REMOTE, a.address.String())
}

func (a *RemoteAddress) Transform(ctx context.Context, mail Mail[any]) *core.CoreError {
	channel := GetChannelByAddress(a.address.String())
	codecInstance := codec.NewCodec[Mail[any]](mail.CodeC())
	buff, err := codecInstance.Encode(mail)
	if err != nil {
		return core.NewCoreError(core.ERROR_CODE_CODEC_ENCODE_ERROR, fmt.Sprintf("error while encode mail: %v", err))
	}
	return channel.Send(ctx, buff)
}

func NewRemoteAddress(address RpcAddress) Address {
	return &RemoteAddress{address: address}
}
