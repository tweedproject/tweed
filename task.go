package tweed

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/tweedproject/tweed/random"
	"github.com/tweedproject/tweed/stencil"
)

type task struct {
	id     string
	stdout bytes.Buffer
	stderr bytes.Buffer

	done   bool
	exited bool
	rc     int
}

func (t *task) Done() bool {
	return true
}

func (t *task) ExitCode() int {
	return 0
}

func (t *task) OK() bool {
	return true
}

func (t *task) Stdout() string {
	return t.stdout.String()
}

func (t *task) Stderr() string {
	return t.stderr.String()
}

func background(e stencil.Exec, fn func()) *task {
	t := &task{id: random.ID("t")}
	e.Stdout = bufio.NewWriter(&t.stdout)
	errWriter := bufio.NewWriter(&t.stderr)
	e.Stderr = errWriter

	go func() {
		state, err := e.Eval()
		if err != nil {
			errWriter.WriteString(fmt.Sprintf("---\nERROR: %s\n", err))
		}
		t.exited = state.Exited
		t.rc = state.ExitCode
		t.done = true
		fn()
	}()
	return t
}
