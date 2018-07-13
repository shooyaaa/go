package uuid

type ID int64

type UUID interface {
	NewUUID() ID
}
