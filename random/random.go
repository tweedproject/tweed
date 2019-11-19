package random

import (
	"crypto/rand"
	"fmt"
)

func panicf(f string, args ...interface{}) {
	panic(fmt.Sprintf(f, args...))
}

func ID(prefix string) string {
	b := make([]byte, 7)
	n, err := rand.Read(b)
	if n != len(b) {
		panicf("unable to generate random id: only read %d bytes from entropy source (wanted 24)", n)
	}
	if err != nil {
		panicf("unable to generate random id: %s", err)
	}

	return fmt.Sprintf("%s-%x", prefix, b)
}
