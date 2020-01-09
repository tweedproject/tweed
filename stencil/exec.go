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
	Args    []string
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
		return []byte{}, fmt.Errorf("%s\nSTDERR: %s\nSTDOUT: %s",
			err, stderr.String(), stdout.String())
	}
	return stdout.Bytes(), nil
}
