package actor

import (
	"context"

	"github.com/shooyaaa/core/uuid"
	"github.com/shooyaaa/log"
)

type ActorType string

type ActorProcessFn[T any] func(T)

type Actor[T Mail[any], D any] interface {
	ActorStartImpl[T]
	ActorStopImpl
	ActorMailboxImpl[T]
	Data() D
	ID() uuid.UUID
}

type ActorStartImpl[T Mail[any]] interface {
	Start(ActorProcessFn[T])
}
type ActorStopImpl interface {
	Stop()
}

type ActorMailboxImpl[T Mail[any]] interface {
	Mailbox() Mailbox[T]
}

func NewActor[T Mail[any], D any](mailboxType MailboxType, id uuid.UUID, data D) Actor[T, D] {
	return &actorImpl[T, D]{
		id:      id,
		mailbox: NewMailbox(mailboxType).(Mailbox[T]),
		data:    data,
	}
}
func (a *actorImpl[T, D]) Stop() {
	a.mailbox.Close(context.Background())
}

type actorImpl[T Mail[any], D any] struct {
	id      uuid.UUID
	mailbox Mailbox[T]
	data    D
}

func (a *actorImpl[T, D]) Mailbox() Mailbox[T] {
	return a.mailbox
}

func (a *actorImpl[T, D]) Data() D {
	return a.data
}

func (a *actorImpl[T, D]) Start(process ActorProcessFn[T]) {
	go func() {
		for {
			msg, err := a.mailbox.Receive(context.Background())
			if err != nil {
				log.ErrorF("error while receive message: %v", err)
				continue
			}
			process(msg)
		}
	}()
}

func (a *actorImpl[T, D]) ID() uuid.UUID {
	return a.id
}
