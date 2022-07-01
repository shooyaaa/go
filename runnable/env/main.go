package env

import (
	"github.com/shooyaaa/config"
	"github.com/shooyaaa/core/library"
	"os"
)

type Env struct {
}

func (c Env) Run() library.Module {
	return &Env{}
}

func (e Env) Start() {
	if _, err := os.Stat(config.TmpDir); err != nil {
		os.Mkdir(config.TmpDir, 0660)
	}
}
func (e Env) Stop() {

}
