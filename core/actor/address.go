package actor

import (
	"context"
	"fmt"

	"github.com/shooyaaa/core"
	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/uuid"
)

type AddressType string

const AddressType_LOCAL AddressType = "local"
const AddressType_REMOTE AddressType = "remote"

type Address interface {
	String() string
	ID() uuid.UUID
	Transfer(ctx context.Context, mail Mail[any]) *core.CoreError
}

type LocalPostManAddress struct {
	postman Postman
}

func (a *LocalPostManAddress) String() string {
	return fmt.Sprintf("%s:%s", AddressType_LOCAL, a.postman.ID())
}

func (a *LocalPostManAddress) ID() uuid.UUID {
	return a.postman.ID()
}

func (a *LocalPostManAddress) Transfer(ctx context.Context, mail Mail[any]) *core.CoreError {
	return a.postman.Receive(ctx, mail)
}

type LocalPostOfficeAddress struct {
	postoffice Postoffice
}

func (a *LocalPostOfficeAddress) String() string {
	return fmt.Sprintf("%s:%s", AddressType_LOCAL, a.postoffice.ID())
}

func (a *LocalPostOfficeAddress) ID() uuid.UUID {
	return a.postoffice.ID()
}

type RemoteAddress struct {
	address RpcAddress
}

func (a *RemoteAddress) String() string {
	return fmt.Sprintf("%s:%s", AddressType_REMOTE, a.address.String())
}

func (a *RemoteAddress) Transfer(ctx context.Context, mail Mail[any]) *core.CoreError {
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
