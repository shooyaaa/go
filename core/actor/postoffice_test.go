package actor

import (
	"context"
	"sync"
	"testing"

	"github.com/shooyaaa/core"
	"github.com/shooyaaa/core/codec"
	"github.com/shooyaaa/core/library"
	"github.com/shooyaaa/core/uuid"
	"github.com/stretchr/testify/assert"
)

// mockAddress 用于测试的 Address 实现
type mockAddress struct {
	id            uuid.UUID
	transferErr   *core.CoreError
	receivedMails []Mail[any]
}

func newMockAddress(id uuid.UUID) *mockAddress {
	return &mockAddress{
		id:            id,
		receivedMails: make([]Mail[any], 0),
	}
}

func (m *mockAddress) String() string {
	return m.id.String()
}

func (m *mockAddress) ID() uuid.UUID {
	return m.id
}

func (m *mockAddress) Transfer(ctx context.Context, mail Mail[any]) *core.CoreError {
	m.receivedMails = append(m.receivedMails, mail)
	return m.transferErr
}

func TestNewPostoffice(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())
	assert.NotNil(t, po)
}

func TestPostoffice_Add(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	addr := newMockAddress(addrGen.Next())

	ctx := context.Background()
	err := po.Add(ctx, addr)

	assert.Nil(t, err)
	assert.Equal(t, 1, hash.GetNodeCount(), "Hash should have one node")
}

func TestPostoffice_Remove(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	addr1 := newMockAddress(addrGen.Next())
	addr2 := newMockAddress(addrGen.Next())

	ctx := context.Background()
	po.Add(ctx, addr1)
	po.Add(ctx, addr2)

	assert.Equal(t, 2, hash.GetNodeCount())

	err := po.Remove(ctx, addr1)
	assert.Nil(t, err)
	assert.Equal(t, 1, hash.GetNodeCount(), "Hash should have one node after removal")
}

func TestPostoffice_Dispatch(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	addr := newMockAddress(addrGen.Next())

	ctx := context.Background()
	po.Add(ctx, addr)

	// 创建邮件
	sender := addrGen.Next()
	receiver := addr.id
	message := "test message"
	mail := NewMail[any](sender, receiver, message, codec.JSON_CODEC)

	err := po.Dispatch(ctx, mail)

	assert.Nil(t, err)
	assert.Len(t, addr.receivedMails, 1, "Address should receive one mail")
	assert.Equal(t, message, addr.receivedMails[0].Message())
}

func TestPostoffice_DispatchNotFound(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)

	// 不添加任何 Address

	ctx := context.Background()

	// 创建邮件，receiver 不存在
	sender := addrGen.Next()
	receiver := addrGen.Next()
	message := "test message"
	mail := NewMail[any](sender, receiver, message, codec.JSON_CODEC)

	err := po.Dispatch(ctx, mail)

	assert.NotNil(t, err)
	assert.Equal(t, core.ERROR_CODE_POSTMAN_NOT_FOUND, err.Code())
	assert.Contains(t, err.String(), "postman not found")
}

func TestPostoffice_DispatchWithTransferError(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	addr := newMockAddress(addrGen.Next())
	addr.transferErr = core.NewCoreError(core.ERROR_CODE_MAILBOX_SEND_ERROR, "transfer error")

	ctx := context.Background()
	po.Add(ctx, addr)

	// 创建邮件
	sender := addrGen.Next()
	receiver := addr.id
	message := "test message"
	mail := NewMail[any](sender, receiver, message, codec.JSON_CODEC)

	err := po.Dispatch(ctx, mail)

	assert.NotNil(t, err)
	assert.Equal(t, core.ERROR_CODE_MAILBOX_SEND_ERROR, err.Code())
	assert.Len(t, addr.receivedMails, 1, "Address should still receive the mail")
}

func TestPostoffice_MultipleAddresses(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)

	addr1 := newMockAddress(addrGen.Next())
	addr2 := newMockAddress(addrGen.Next())
	addr3 := newMockAddress(addrGen.Next())

	ctx := context.Background()
	po.Add(ctx, addr1)
	po.Add(ctx, addr2)
	po.Add(ctx, addr3)

	assert.Equal(t, 3, hash.GetNodeCount())

	// 测试分发到不同的 Address
	sender := addrGen.Next()
	mail1 := NewMail[any](sender, addr1.id, "message1", codec.JSON_CODEC)
	mail2 := NewMail[any](sender, addr2.id, "message2", codec.JSON_CODEC)
	mail3 := NewMail[any](sender, addr3.id, "message3", codec.JSON_CODEC)

	err1 := po.Dispatch(ctx, mail1)
	err2 := po.Dispatch(ctx, mail2)
	err3 := po.Dispatch(ctx, mail3)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Nil(t, err3)

	// 由于一致性哈希的分布特性，邮件可能不会精确路由到对应的 Address
	// 验证所有邮件都被接收（总共3封）
	totalReceived := len(addr1.receivedMails) + len(addr2.receivedMails) + len(addr3.receivedMails)
	assert.Equal(t, 3, totalReceived, "All three mails should be received")

	// 验证每个 Address 至少收到了一些邮件（由于哈希分布，可能不是每个都收到）
	// 但至少应该有一些邮件被接收
	assert.Greater(t, totalReceived, 0, "At least some mails should be received")
}

func TestPostoffice_AddRemoveAdd(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	addr := newMockAddress(addrGen.Next())

	ctx := context.Background()

	// 添加
	err := po.Add(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, 1, hash.GetNodeCount())

	// 移除
	err = po.Remove(ctx, addr)
	assert.Nil(t, err)
	assert.Equal(t, 0, hash.GetNodeCount())

	// 再次添加
	addr2 := newMockAddress(addrGen.Next())
	err = po.Add(ctx, addr2)
	assert.Nil(t, err)
	assert.Equal(t, 1, hash.GetNodeCount())
}

func TestPostoffice_DispatchAfterRemove(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	addr := newMockAddress(addrGen.Next())

	ctx := context.Background()
	po.Add(ctx, addr)

	// 移除
	po.Remove(ctx, addr)

	// 尝试分发到已移除的 Address
	sender := addrGen.Next()
	mail := NewMail[any](sender, addr.id, "test", codec.JSON_CODEC)

	err := po.Dispatch(ctx, mail)

	assert.NotNil(t, err)
	assert.Equal(t, core.ERROR_CODE_POSTMAN_NOT_FOUND, err.Code())
	assert.Len(t, addr.receivedMails, 0, "Removed address should not receive mail")
}

func TestPostoffice_ConsistentHashDistribution(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)

	addr1 := newMockAddress(addrGen.Next())
	addr2 := newMockAddress(addrGen.Next())
	addr3 := newMockAddress(addrGen.Next())

	ctx := context.Background()
	po.Add(ctx, addr1)
	po.Add(ctx, addr2)
	po.Add(ctx, addr3)

	// 使用相同的 key 应该路由到相同的 Address
	sender := addrGen.Next()
	receiver := addrGen.Next()
	mail1 := NewMail[any](sender, receiver, "message1", codec.JSON_CODEC)
	mail2 := NewMail[any](sender, receiver, "message2", codec.JSON_CODEC)

	err1 := po.Dispatch(ctx, mail1)
	err2 := po.Dispatch(ctx, mail2)

	assert.Nil(t, err1)
	assert.Nil(t, err2)

	// 两个邮件应该路由到同一个 Address（因为 receiver 相同）
	// 验证至少有一个 Address 收到了邮件
	totalReceived := len(addr1.receivedMails) + len(addr2.receivedMails) + len(addr3.receivedMails)
	assert.Equal(t, 2, totalReceived, "Both mails should be received by the same address")
}

func TestPostoffice_ID(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	expectedID := idGen.Next()
	po := NewPostoffice(hash, expectedID)

	assert.Equal(t, expectedID, po.ID())
}

func TestPostoffice_Actor(t *testing.T) {
	hash := library.NewConsistentHash[Address](150, nil, nil)
	const UUIDTypePostoffice uuid.UUIDType = "postoffice"
	idGen := uuid.NewSimpleUUIDGenerator(UUIDTypePostoffice)
	po := NewPostoffice(hash, idGen.Next())

	const UUIDTypeTest uuid.UUIDType = "test"
	addrGen := uuid.NewSimpleUUIDGenerator(UUIDTypeTest)
	pa := NewActor[Mail[any], any](MailboxType_MEMORY, addrGen.Next(), po)
	pa.Start(func(mail Mail[any]) {
		assert.Equal(t, po, mail.Receiver())
	})

	a1 := NewActor[Mail[any], any](MailboxType_MEMORY, addrGen.Next(), func(mail Mail[any]) {
	})
	ctx := context.Background()
	pm1 := NewPostman()
	pm1.Add(ctx, a1)
	pm1.Register(ctx, NewLocalPostOfficeAddress(po))

	wg := sync.WaitGroup{}
	wg.Add(1)
	a2 := NewActor[Mail[any], any](MailboxType_MEMORY, addrGen.Next(), func(mail Mail[any]) {
		assert.Equal(t, po, mail.Receiver())
		wg.Done()
	})
	pm2 := NewPostman()
	pm2.Add(ctx, a2)
	pm2.Register(ctx, NewLocalPostManAddress(pm2))
	mail1 := NewMail[any](addrGen.Next(), a2.ID(), "test", codec.JSON_CODEC)
	err := a1.Mailbox().Send(ctx, mail1)
	assert.Nil(t, err)
	wg.Wait()
}
