package actor

import (
	"context"
	"fmt"
	"sync"

	"github.com/shooyaaa/core"
	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/uuid"
	"github.com/shooyaaa/log"
)

type PostManAddressType string

const PostManAddressType_LOCAL PostManAddressType = "local"
const PostManAddressType_REMOTE PostManAddressType = "remote"

type PostManAddress interface {
	String() string
	Transform(ctx context.Context, mail Mail[any]) *core.CoreError
}

type LocalPostManAddress struct {
	postman Postman
}

func (a *LocalPostManAddress) String() string {
	return fmt.Sprintf("%s:%s", PostManAddressType_LOCAL, a.postman.ID())
}

func (a *LocalPostManAddress) Transform(ctx context.Context, mail Mail[any]) *core.CoreError {
	return a.postman.Receive(ctx, mail)
}
func NewLocalPostManAddress(postman Postman) PostManAddress {
	return &LocalPostManAddress{postman: postman}
}

type RemotePostManAddress struct {
	address RpcAddress
}

func (a *RemotePostManAddress) String() string {
	return fmt.Sprintf("%s:%s", PostManAddressType_REMOTE, a.address.String())
}

func (a *RemotePostManAddress) Transform(ctx context.Context, mail Mail[any]) *core.CoreError {
	channel := GetChannelByAddress(a.address.String())
	codecInstance := codec.NewCodec[Mail[any]](mail.CodeC())
	buff, err := codecInstance.Encode(mail)
	if err != nil {
		return core.NewCoreError(core.ERROR_CODE_CODEC_ENCODE_ERROR, fmt.Sprintf("error while encode mail: %v", err))
	}
	return channel.Send(ctx, buff)
}

func NewRemotePostManAddress(address RpcAddress) PostManAddress {
	return &RemotePostManAddress{address: address}
}

type Postman interface {
	Add(ctx context.Context, a Actor[Mail[any], any]) *core.CoreError
	Deliver(ctx context.Context, mail Mail[any]) *core.CoreError
	Remove(ctx context.Context, id uuid.UUID) *core.CoreError
	Dispatch(ctx context.Context, mail Mail[any]) *core.CoreError
	Receive(ctx context.Context, mail Mail[any]) *core.CoreError
	ID() uuid.UUID
}

type postmanImpl struct {
	actors sync.Map
	id     uuid.UUID
}

func NewPostman() Postman {
	return &postmanImpl{}
}

func (m *postmanImpl) ID() uuid.UUID {
	return m.id
}

func (m *postmanImpl) Add(ctx context.Context, a Actor[Mail[any], any]) *core.CoreError {
	m.actors.Store(a.ID(), a)
	a.Mailbox().Gather(func(mail Mail[any]) {
		err2 := m.Deliver(ctx, mail)
		if err2 != nil {
			log.ErrorF("error while deliver message: %s", err2.String())
		}
	})
	return nil
}

func (m *postmanImpl) Receive(ctx context.Context, mail Mail[any]) *core.CoreError {
	a, ok := m.actors.Load(mail.Receiver())
	if ok {
		err1 := a.(Actor[Mail[any], any]).Mailbox().Send(ctx, mail)
		if err1 != nil {
			return core.NewCoreError(core.ERROR_CODE_MAILBOX_SEND_ERROR, err1.Error())
		}
	} else {
		receiver := mail.Receiver()
		return core.NewCoreError(core.ERROR_CODE_ACTOR_NOT_FOUND, fmt.Sprintf("postman receive a mail but actor not found: %s", (&receiver).String()))
	}
	return nil
}

func (m *postmanImpl) Deliver(ctx context.Context, mail Mail[any]) *core.CoreError {
	a, ok := m.actors.Load(mail.Receiver())
	if ok {
		err1 := a.(Actor[Mail[any], any]).Mailbox().Send(ctx, mail)
		if err1 != nil {
			return core.NewCoreError(core.ERROR_CODE_MAILBOX_SEND_ERROR, err1.Error())
		}
	} else {
		return m.Dispatch(ctx, mail)
	}
	return nil

}

func (m *postmanImpl) Dispatch(ctx context.Context, mail Mail[any]) *core.CoreError {
	if m.postoffice == nil {
		return core.NewCoreError(core.ERROR_CODE_POSTOFFICE_NOT_REGISTERED, "postoffice not registered")
	}
	return m.postoffice.Dispatch(ctx, mail)
}

func (m *postmanImpl) Remove(ctx context.Context, id uuid.UUID) *core.CoreError {
	_, ok := m.actors.LoadAndDelete(id)
	if !ok {
		return core.NewCoreError(core.ERROR_CODE_ACTOR_NOT_FOUND, fmt.Sprintf("actor not found: %s", (&id).String()))
	}
	return nil
}
