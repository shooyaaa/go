package actor

import (
	"context"
	"testing"

	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/uuid"
	"github.com/stretchr/testify/assert"
)

func TestActor(t *testing.T) {
	const UUIDTypeActor uuid.UUIDType = "actor"
	const UUIDTypeMail uuid.UUIDType = "mail"
	actor := NewActor[Mail[any], string](MailboxType_MEMORY, uuid.NewSimpleUUIDGenerator(UUIDTypeActor).Next(), "test")
	actor.Start(func(msg Mail[any]) {
		assert.Equal(t, any("test"), msg.Message())
		assert.Equal(t, 1, actor.Data())
		actor.Stop()
	})
	err := actor.Mailbox().Send(context.Background(), NewMail[any](uuid.NewSimpleUUIDGenerator(UUIDTypeActor).Next(), uuid.NewSimpleUUIDGenerator(UUIDTypeMail).Next(), any("test"), codec.JSON_CODEC))
	assert.NoError(t, err)
	actor.Stop()
}
