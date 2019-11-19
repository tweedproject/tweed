package tweed

import (
	"fmt"
	"time"
)

type oops struct {
	handler string
	remote  string
	request []byte
	dated   time.Time
	err     error
}

func (o oops) message() string {
	return o.err.Error()
}

type Errors struct {
	all []error
}

func (e Errors) Error() string {
	var s string
	for _, err := range e.all {
		s += fmt.Sprintf(" - %s\n", err)
	}
	return s
}

func (c *Core) KeepErrors(n int) {
	c.oops.max = n
}
