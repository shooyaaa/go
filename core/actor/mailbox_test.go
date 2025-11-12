package actor

import (
	"context"
	"testing"

	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMemoryMailbox(t *testing.T) {
	mb := NewMemoryMailbox[Mail[any]]("test")

	const UUIDTypeActor uuid.UUIDType = "actor"
	// 创建测试消息
	sender := uuid.NewSimpleUUIDGenerator(UUIDTypeActor).Next()
	receiver := uuid.NewSimpleUUIDGenerator(UUIDTypeActor).Next()
	message := any("test")
	mail := NewMail[any](sender, receiver, message, codec.JSON_CODEC)
	err := mb.Send(context.Background(), mail)
	assert.NoError(t, err)
	data, err := mb.Receive(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, any("test"), data.Message())
	mb.Close(context.Background())
	assert.NoError(t, err)
}
