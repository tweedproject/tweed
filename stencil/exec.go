package stencil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type ProcessState struct {
	ExitCode int
	Exited   bool
}

type Exec struct {
	Run     string
	Env     []string
	Stdout  io.Writer
	Stderr  io.Writer
	Mounts  []string
	Stencil *Stencil
}

func Run(e Exec) ([]byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	e.Stdout = bufio.NewWriter(&stdout)
	e.Stderr = bufio.NewWriter(&stderr)
	_, err := e.Eval()
	if err != nil {
		return []byte(fmt.Sprintf("--STDOUT--\n%s\n---STDERR--\n%s",
			stdout.String(), stderr.String())), err
	}
	return stdout.Bytes(), nil
}
