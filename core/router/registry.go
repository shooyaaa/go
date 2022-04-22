package router

import (
	"errors"
	"github.com/shooyaaa/core"
)

type DummyRegistry struct {
	tables map[string]chan *Package
}

func (dr *DummyRegistry) Get(id string) (chan *Package, error) {
	entry, ok := dr.tables[id]
	if ok {
		return entry, nil
	}
	return nil, errors.New(core.NOT_FOUND)
}

func (dr *DummyRegistry) Set(id string, p chan *Package) error {
	dr.tables[id] = p
	return nil
}
func (dr *DummyRegistry) Delete(name string) {
	delete(dr.tables, name)
}

type TcpRegistry struct {
	tables map[string]string
}

func (tr *TcpRegistry) Get(id string) (string, error) {
	entry, ok := tr.tables[id]
	if ok {
		return entry, nil
	}
	return "", errors.New(core.NOT_FOUND)
}

func (tr *TcpRegistry) Set(id string, value string) error {
	tr.tables[id] = value
	return nil
}

func (tr *TcpRegistry) Remove(name string) {
	delete(tr.tables, name)
}
