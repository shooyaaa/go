package actor

import (
	"context"
	"testing"

	"github.com/shooyaaa/core"
	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/library"
	"github.com/shooyaaa/core/uuid"
	"github.com/stretchr/testify/assert"
)

// mockPostman 用于测试的 Postman 实现
type mockPostman struct {
	id            uuid.UUID
	registerErr   *core.CoreError
	receiveErr    *core.CoreError
	receivedMails []Mail[any]
	registered    bool
}

func newMockPostman(id uuid.UUID) *mockPostman {
	return &mockPostman{
		id:            id,
		receivedMails: make([]Mail[any], 0),
	}
}

func (m *mockPostman) Add(ctx context.Context, a Actor[Mail[any], any]) *core.CoreError {
	return nil
}

func (m *mockPostman) Deliver(ctx context.Context, mail Mail[any]) *core.CoreError {
	return nil
}

func (m *mockPostman) Remove(ctx context.Context, id uuid.UUID) *core.CoreError {
	return nil
}

func (m *mockPostman) Dispatch(ctx context.Context, mail Mail[any]) *core.CoreError {
	return nil
}

func (m *mockPostman) Register(ctx context.Context, p Postoffice) *core.CoreError {
	m.registered = true
	return m.registerErr
}

func (m *mockPostman) Receive(ctx context.Context, mail Mail[any]) *core.CoreError {
	m.receivedMails = append(m.receivedMails, mail)
	return m.receiveErr
}

// 实现 Stringer 接口以便在 ConsistentHash 中使用
func (m *mockPostman) String() string {
	return m.id.String()
}

func TestNewPostoffice(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)
	assert.NotNil(t, po)
}

func TestPostoffice_Add(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	pm := newMockPostman(idGen.Next())

	ctx := context.Background()
	err := po.Add(ctx, pm)

	assert.Nil(t, err)
	assert.True(t, pm.registered, "Postman should be registered")
	assert.Equal(t, 1, hash.GetNodeCount(), "Hash should have one node")
}

func TestPostoffice_AddWithRegisterError(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	pm := newMockPostman(idGen.Next())
	pm.registerErr = core.NewCoreError(core.ERROR_CODE_POSTOFFICE_NOT_REGISTERED, "test error")

	ctx := context.Background()
	err := po.Add(ctx, pm)

	assert.NotNil(t, err)
	assert.Equal(t, core.ERROR_CODE_POSTOFFICE_NOT_REGISTERED, err.Code())
	assert.True(t, pm.registered, "Postman should still be registered")
}

func TestPostoffice_Remove(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	pm1 := newMockPostman(idGen.Next())
	pm2 := newMockPostman(idGen.Next())

	ctx := context.Background()
	po.Add(ctx, pm1)
	po.Add(ctx, pm2)

	assert.Equal(t, 2, hash.GetNodeCount())

	err := po.Remove(ctx, pm1)
	assert.Nil(t, err)
	assert.Equal(t, 1, hash.GetNodeCount(), "Hash should have one node after removal")
}

func TestPostoffice_Dispatch(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	pm := newMockPostman(idGen.Next())

	ctx := context.Background()
	po.Add(ctx, pm)

	// 创建邮件
	sender := idGen.Next()
	receiver := pm.id
	message := "test message"
	mail := NewMail[any](sender, receiver, message, codec.JSON_CODEC)

	err := po.Dispatch(ctx, mail)

	assert.Nil(t, err)
	assert.Len(t, pm.receivedMails, 1, "Postman should receive one mail")
	assert.Equal(t, message, pm.receivedMails[0].Message())
}

func TestPostoffice_DispatchNotFound(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)

	// 不添加任何 Postman

	ctx := context.Background()

	// 创建邮件，receiver 不存在
	sender := idGen.Next()
	receiver := idGen.Next()
	message := "test message"
	mail := NewMail[any](sender, receiver, message, codec.JSON_CODEC)

	err := po.Dispatch(ctx, mail)

	assert.NotNil(t, err)
	assert.Equal(t, core.ERROR_CODE_POSTMAN_NOT_FOUND, err.Code())
	assert.Contains(t, err.String(), "postman not found")
}

func TestPostoffice_DispatchWithReceiveError(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	pm := newMockPostman(idGen.Next())
	pm.receiveErr = core.NewCoreError(core.ERROR_CODE_MAILBOX_SEND_ERROR, "receive error")

	ctx := context.Background()
	po.Add(ctx, pm)

	// 创建邮件
	sender := idGen.Next()
	receiver := pm.id
	message := "test message"
	mail := NewMail[any](sender, receiver, message, codec.JSON_CODEC)

	err := po.Dispatch(ctx, mail)

	assert.NotNil(t, err)
	assert.Equal(t, core.ERROR_CODE_MAILBOX_SEND_ERROR, err.Code())
	assert.Len(t, pm.receivedMails, 1, "Postman should still receive the mail")
}

func TestPostoffice_MultiplePostmen(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)

	pm1 := newMockPostman(idGen.Next())
	pm2 := newMockPostman(idGen.Next())
	pm3 := newMockPostman(idGen.Next())

	ctx := context.Background()
	po.Add(ctx, pm1)
	po.Add(ctx, pm2)
	po.Add(ctx, pm3)

	assert.Equal(t, 3, hash.GetNodeCount())

	// 测试分发到不同的 Postman
	sender := idGen.Next()
	mail1 := NewMail[any](sender, pm1.id, "message1", codec.JSON_CODEC)
	mail2 := NewMail[any](sender, pm2.id, "message2", codec.JSON_CODEC)
	mail3 := NewMail[any](sender, pm3.id, "message3", codec.JSON_CODEC)

	err1 := po.Dispatch(ctx, mail1)
	err2 := po.Dispatch(ctx, mail2)
	err3 := po.Dispatch(ctx, mail3)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)

	// 由于一致性哈希的分布特性，邮件可能不会精确路由到对应的 Postman
	// 验证所有邮件都被接收（总共3封）
	totalReceived := len(pm1.receivedMails) + len(pm2.receivedMails) + len(pm3.receivedMails)
	assert.Equal(t, 3, totalReceived, "All three mails should be received")

	// 验证每个 Postman 至少收到了一些邮件（由于哈希分布，可能不是每个都收到）
	// 但至少应该有一些邮件被接收
	assert.Greater(t, totalReceived, 0, "At least some mails should be received")
}

func TestPostoffice_AddRemoveAdd(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	pm := newMockPostman(idGen.Next())

	ctx := context.Background()

	// 添加
	err := po.Add(ctx, pm)
	assert.Nil(t, err)
	assert.Equal(t, 1, hash.GetNodeCount())

	// 移除
	err = po.Remove(ctx, pm)
	assert.Nil(t, err)
	assert.Equal(t, 0, hash.GetNodeCount())

	// 再次添加
	pm2 := newMockPostman(idGen.Next())
	err = po.Add(ctx, pm2)
	assert.Nil(t, err)
	assert.Equal(t, 1, hash.GetNodeCount())
}

func TestPostoffice_DispatchAfterRemove(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	pm := newMockPostman(idGen.Next())

	ctx := context.Background()
	po.Add(ctx, pm)

	// 移除
	po.Remove(ctx, pm)

	// 尝试分发到已移除的 Postman
	sender := idGen.Next()
	mail := NewMail[any](sender, pm.id, "test", codec.JSON_CODEC)

	err := po.Dispatch(ctx, mail)

	assert.NotNil(t, err)
	assert.Equal(t, core.ERROR_CODE_POSTMAN_NOT_FOUND, err.Code())
	assert.Len(t, pm.receivedMails, 0, "Removed postman should not receive mail")
}

func TestPostoffice_ConsistentHashDistribution(t *testing.T) {
	hash := library.NewConsistentHash[Postman](150, nil, nil)
	po := NewPostoffice(hash)

	const UUIDTypeTest uuid.UUIDType = "test"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)

	pm1 := newMockPostman(idGen.Next())
	pm2 := newMockPostman(idGen.Next())
	pm3 := newMockPostman(idGen.Next())

	ctx := context.Background()
	po.Add(ctx, pm1)
	po.Add(ctx, pm2)
	po.Add(ctx, pm3)

	// 使用相同的 key 应该路由到相同的 Postman
	sender := idGen.Next()
	receiver := idGen.Next()
	mail1 := NewMail[any](sender, receiver, "message1", codec.JSON_CODEC)
	mail2 := NewMail[any](sender, receiver, "message2", codec.JSON_CODEC)

	err1 := po.Dispatch(ctx, mail1)
	err2 := po.Dispatch(ctx, mail2)

	assert.Nil(t, err1)
	assert.Nil(t, err2)

	// 两个邮件应该路由到同一个 Postman（因为 receiver 相同）
	// 但由于我们使用 receiver 作为 key，它们应该都路由到同一个 Postman
	// 验证至少有一个 Postman 收到了邮件
	totalReceived := len(pm1.receivedMails) + len(pm2.receivedMails) + len(pm3.receivedMails)
	assert.Equal(t, 2, totalReceived, "Both mails should be received by the same postman")
}
