package gate

import (
	"errors"
	"github.com/shooyaaa/core/network"
	"github.com/shooyaaa/core/router"
	"github.com/shooyaaa/core/session"
	"github.com/shooyaaa/log"
)

type Gate struct {
	network.Server
	list     map[int64]*session.Session
	WaitChan chan *session.Session
}

func NewGate() *Gate {
	gate := Gate{
		Server: &network.Tcp{
			Id:        nil,
			HeartBeat: 0,
		},
		list:     make(map[int64]*session.Session),
		WaitChan: make(chan *session.Session, 0),
	}
	return &gate
}

func (gate *Gate) Run() {
	log.Debug("gate start Work")
	go gate.Accept()
}

func (gate *Gate) Accept() {
	for {
		select {
		case session := <-gate.WaitChan:
			gate.list[session.Id] = session
			log.DebugF("New session %v", session.Id)
		}
	}
}

func (gate *Gate) Get(id int64) (*session.Session, error) {
	sessionItem, ok := gate.list[id]
	if !ok {
		return sessionItem, errors.New("Invalid Session id %d")
	}
	return sessionItem, nil
}

func (gate *Gate) PackageHandler(p *router.Package) {
}
func (gate *Gate) SessionClose(id int64) {
	log.DebugF("session close in gate: %v", id)
	delete(gate.list, id)
}
