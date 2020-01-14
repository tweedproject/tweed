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
	Cwd     string
	Stdout  io.Writer
	Stderr  io.Writer
	Mounts  []Mount
	Stencil *Stencil
}

type Mount struct {
	Source      string
	Destination string
	Writable    bool
}

func Run(e Exec) ([]byte, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	e.Stdout = bufio.NewWriter(&stdout)
	e.Stderr = bufio.NewWriter(&stderr)
	state, err := e.Eval()
	if err != nil {
		return nil, fmt.Errorf("failed to run procces: %s\n%s", err, stderr.Bytes())
	}
	if state.ExitCode != 0 {
		return nil, fmt.Errorf("process exited with: %d\n%s", state.ExitCode, stderr.Bytes())
	}
	return stdout.Bytes(), nil
}
