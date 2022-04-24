package storage

type Cache interface {
	GetInt(string) (int, error)
	GetInt64(string) (int64, error)
	GetString(string) string
	SetInt(string, int) error
	SetInt64(string, int64) error
	SetString(string, string) error
	Delete(string) error
	Init(map[string]interface{})
}
