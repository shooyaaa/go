package types

type Simple struct {
	Counter int64
}

func (s *Simple) NewUUID() int64 {
	s.Counter++
	return s.Counter
}
