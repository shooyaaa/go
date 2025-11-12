package uuid

import "fmt"

type UUIDType string
type UUID struct {
	ID   int64
	Type UUIDType
}

func (u *UUID) String() string {
	return fmt.Sprintf("%s:%d", u.Type, u.ID)
}

type SimpleUUIDGenerator interface {
	Next() UUID
}

type SimpleUUIDGeneratorImpl struct {
	Type UUIDType
	ID   int64
}

func (u *SimpleUUIDGeneratorImpl) Next() UUID {
	u.ID++
	return UUID{Type: u.Type, ID: u.ID}
}

func NewSimpleUUIDGenerator(t UUIDType) SimpleUUIDGenerator {
	return &SimpleUUIDGeneratorImpl{Type: t, ID: 0}
}
