package runnable

import (
	gate2 "github.com/shooyaaa/runnable/gate"
	"testing"
)

func TestGate(t *testing.T) {
	gate := gate2.NewGate()
	gate.Listen("127.0.0.1:9797")
}
