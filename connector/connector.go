package connector

import (
	"github.com/shooyaaa/manager"
)

type Connector interface {
	Run() manager.Session 
}