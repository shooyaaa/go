package router

import (
	"errors"

	"github.com/shooyaaa/core/storage"
)

type Registry interface {
	Get(id string) interface{}
}

type DummyRegistry struct {
	tables map[string]chan *Package
}

func (dr *DummyRegistry) Get(id string) (chan *Package, error) {
	entry, ok := dr.tables[id]
	if ok {
		return entry, nil
	}
	return nil, errors.New("not found")
}

func (dr *DummyRegistry) Set(id string, p chan *Package) error {
	dr.tables[id] = p
	return nil
}
func (dr *DummyRegistry) Delete(name string) {
	delete(dr.tables, name)
}

type TcpRegistry struct {
	tables storage.Cache
}

func (tr *TcpRegistry) Get(id string) (string, error) {
	if tr.tables == nil {
		return "", errors.New("not found")
	}
	entry := tr.tables.GetString(id)
	if len(entry) > 0 {
		return entry, nil
	}
	return "", errors.New("not found")
}

func (tr *TcpRegistry) Set(id string, value string) error {
	tr.tables.SetString(id, value)
	return nil
}

func (tr *TcpRegistry) Delete(name string) {
	tr.tables.Delete(name)
}

func (tr *TcpRegistry) Init(params map[string]interface{}) {
	tr.tables.Init(params)
}

func (tr *TcpRegistry) CentralAddress() string {
	return "central_address.local"
}
