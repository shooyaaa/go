package actor

import (
	"context"
	"fmt"

	"github.com/shooyaaa/core/uuid"
)

type MailboxType string

const MailboxType_MEMORY MailboxType = "memory"
const MailboxType_TCP MailboxType = "tcp"
const MailboxType_UDP MailboxType = "udp"
const MailboxType_IPC MailboxType = "ipc"
const MailboxType_UNIX MailboxType = "unix"
const MailboxType_FILE MailboxType = "file"
const MailboxType_PIPE MailboxType = "pipe"

type Mailbox[T Mail[any]] interface {
	MailboxSender[T]
	MailboxReceiver[T]
	MailboxGather[T]
	MailboxCloser[T]
	ID() uuid.UUID
}
type MailboxReceiver[T Mail[any]] interface {
	Receive(ctx context.Context) (T, error)
}
type MailboxSender[T Mail[any]] interface {
	Send(ctx context.Context, data T) error
}
type MailboxGather[T Mail[any]] interface {
	Gather(func(T))
}
type MailboxCloser[T Mail[any]] interface {
	Close(ctx context.Context) error
}

type memoryMailbox[T Mail[any]] struct {
	name   string
	recvCh chan T
	sendCh chan T
	id     uuid.UUID
}

func NewMemoryMailbox[T Mail[any]](name string) Mailbox[T] {
	return &memoryMailbox[T]{name: name, recvCh: make(chan T, 100), sendCh: make(chan T, 100)}
}

func (mb *memoryMailbox[T]) Send(ctx context.Context, data T) error {
	mb.sendCh <- data
	return nil
}

func (mb *memoryMailbox[T]) Receive(ctx context.Context) (T, error) {
	data, ok := <-mb.recvCh
	if !ok {
		return data, fmt.Errorf("channel closed")
	}
	return data, nil
}

func (mb *memoryMailbox[T]) Close(ctx context.Context) error {
	close(mb.recvCh)
	close(mb.sendCh)
	return nil
}

func (mb *memoryMailbox[T]) Gather(fn func(T)) {
	go func() {
		for {
			data, ok := <-mb.sendCh
			if !ok {
				break
			}
			fn(data)
		}
	}()
}

func (mb *memoryMailbox[T]) ID() uuid.UUID {
	return mb.id
}

func NewMailbox(mailboxType MailboxType) Mailbox[Mail[any]] {
	switch mailboxType {
	case MailboxType_MEMORY:
		return NewMemoryMailbox[Mail[any]](string(mailboxType))
	default:
		panic(fmt.Sprintf("unknown mailbox type: %v", mailboxType))
	}
}
