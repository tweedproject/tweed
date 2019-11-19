package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func bail(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} %s\n", err)
		os.Exit(2)
	}
}

func GonnaNeedATweed() {
	ok := true
	if opts.Tweed == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--tweed} flag given, and @W{$TWEED_URL} not set.\n")
		ok = false
	}
	if opts.Username == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--username} flag given, and @W{$TWEED_USERNAME} not set.\n")
		ok = false
	}
	if opts.Password == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--password} flag given, and @W{$TWEED_PASSWORD} not set.\n")
		ok = false
	}
	if !ok {
		os.Exit(1)
	}
}
