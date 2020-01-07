package stencil

import (
	"log"
	"runtime"

	"github.com/opencontainers/runc/libcontainer"
)

func RuncInit() {
	runtime.GOMAXPROCS(1)
	runtime.LockOSThread()
	factory, _ := libcontainer.New("")
	if err := factory.StartInitialization(); err != nil {
		log.Fatal(err)
	}
	panic("--this line should have never been executed, congratulations--")
}
