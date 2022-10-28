package session

type Manager struct {
	s map[int64]Session
}

func (m *Manager) Init() {
	m.s = map[int64]Session{}
}

func (m *Manager) Add(s Session) {
	m.s[s.Id] = s
}

func (m *Manager) RemoveSession(s Session) {
	delete(m.s, s.Id)
}

func (m *Manager) RemoveId(id int64) {
	delete(m.s, id)
}

func (m *Manager) MultiCastSession(session []Session, msg Op) {
	ids := []int64{}
	for _, s := range session {
		ids = append(ids, s.Id)
	}
	m.MultiCastID(ids, msg)
}

func (m *Manager) MultiCastID(ids []int64, msg Op) {
	for _, id := range ids {
		if s, ok := m.s[id]; ok {
			s.Write(msg)
		}
	}
}

func (m *Manager) Broadcast(msg Op) {
	ids := []int64{}
	for id, _ := range m.s {
		ids = append(ids, id)
	}
	m.MultiCastID(ids, msg)
}
