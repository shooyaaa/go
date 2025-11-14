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
	id := a.postman.ID()
	return fmt.Sprintf("%s:%s", AddressType_LOCAL, (&id).String())
}

func (a *LocalPostManAddress) ID() uuid.UUID {
	return a.postman.ID()
}

func (a *LocalPostManAddress) Transfer(ctx context.Context, mail Mail[any]) *core.CoreError {
	return a.postman.Receive(ctx, mail)
}

func NewLocalPostManAddress(postman Postman) Address {
	return &LocalPostManAddress{postman: postman}
}

type LocalPostOfficeAddress struct {
	postoffice Postoffice
}

func (a *LocalPostOfficeAddress) String() string {
	id := a.postoffice.ID()
	return fmt.Sprintf("%s:%s", AddressType_LOCAL, (&id).String())
}

func (a *LocalPostOfficeAddress) ID() uuid.UUID {
	return a.postoffice.ID()
}

func (a *LocalPostOfficeAddress) Transfer(ctx context.Context, mail Mail[any]) *core.CoreError {
	return a.postoffice.Dispatch(ctx, mail)
}
func NewLocalPostOfficeAddress(postoffice Postoffice) Address {
	return &LocalPostOfficeAddress{postoffice: postoffice}
}

type RemoteAddress struct {
	address RpcAddress
	id      uuid.UUID
}

func (a *RemoteAddress) String() string {
	return fmt.Sprintf("%s:%s", AddressType_REMOTE, a.address.String())
}

func (a *RemoteAddress) ID() uuid.UUID {
	return a.id
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

func NewRemoteAddress(address RpcAddress, id uuid.UUID) Address {
	return &RemoteAddress{address: address, id: id}
}
