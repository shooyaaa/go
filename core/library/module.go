package library

import (
	"sync"
)

type Module interface {
	Start()
	Stop()
}

type Entry struct {
	module    Module
	closeChan chan int
}

type moduleManager struct {
	closeChan chan int
	wg        sync.WaitGroup
	entries   []Entry
}

func (m *moduleManager) Load(module Module) {
	entry := Entry{
		module:    module,
		closeChan: make(chan int),
	}
	m.entries = append(m.entries, entry)
}

func (m *moduleManager) Run() {
	m.wg.Add(1)
	for _, entry := range m.entries {
		go func(entry Entry) {
			go entry.module.Start()
			<-entry.closeChan
			entry.module.Stop()
			m.wg.Done()
		}(entry)
	}
	m.wg.Wait()
}

func (m *moduleManager) Exit() {
	for _, entry := range m.entries {
		entry.closeChan <- 1
	}
}

var ModuleManager *moduleManager

func init() {
	ModuleManager = &moduleManager{
		closeChan: make(chan int),
		wg:        sync.WaitGroup{},
		entries:   make([]Entry, 0),
	}
}
