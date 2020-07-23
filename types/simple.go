package types

type Simple struct {
	Counter ID
}

func (s *Simple) NewUUID() ID {
	s.Counter++
	return s.Counter
}
