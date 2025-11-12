package actor

import (
	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/uuid"
)

type Mail[M any] interface {
	Sender() uuid.UUID
	Receiver() uuid.UUID
	Message() M
	CodeC() codec.CODEC_TYPE
}

type mailImpl[M any] struct {
	sender   uuid.UUID
	receiver uuid.UUID
	message  M
	codec    codec.CODEC_TYPE
}

func NewMail[M any](sender uuid.UUID, receiver uuid.UUID, message M, codec codec.CODEC_TYPE) Mail[M] {
	return &mailImpl[M]{sender: sender, receiver: receiver, message: message, codec: codec}
}

func (m *mailImpl[M]) CodeC() codec.CODEC_TYPE {
	return m.codec
}

func (m *mailImpl[M]) Sender() uuid.UUID {
	return m.sender
}

func (m *mailImpl[M]) Receiver() uuid.UUID {
	return m.receiver
}

func (m *mailImpl[M]) Message() M {
	return m.message
}
