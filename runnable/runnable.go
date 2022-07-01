package runnable

import "github.com/shooyaaa/core/library"

type Runnable interface {
	Run() *library.Module
}
