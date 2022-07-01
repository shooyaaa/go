package main

import (
	"github.com/shooyaaa/core/library"
	"github.com/shooyaaa/runnable/cron"
	"os"
	"os/signal"
)

func main() {
	library.ModuleManager.Load(cron.Cron{}.Run())
	go library.ModuleManager.Run()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}
