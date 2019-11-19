package tweed

import (
	"bytes"

	"github.com/tweedproject/tweed/random"
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

func background(e Exec, fn func()) *task {
	e.Stdout = make(chan string, 0)
	e.Stderr = make(chan string, 0)
	e.Done = make(chan int, 1)

	t := &task{id: random.ID("t")}
	go func() {
		for s := range e.Stdout {
			t.stdout.Write([]byte(s))
		}
	}()
	go func() {
		for s := range e.Stderr {
			t.stderr.Write([]byte(s))
		}
	}()
	go func() {
		for rc := range e.Done {
			t.exited = true
			t.rc = rc
		}
		t.done = true
		fn()
	}()

	go e.run()
	return t
}
