package stencil

import "io"

type RunArg struct {
	Run    string
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
	Done   chan int
	Store  string
}
