package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"
)

func bail(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} %s\n", err)
		os.Exit(2)
	}
}

func GonnaNeedATweed() {
	ok := true
	if Tweed.Tweed == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--tweed} flag given, and @W{$TWEED_URL} not set.\n")
		ok = false
	}
	if Tweed.Username == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--username} flag given, and @W{$TWEED_USERNAME} not set.\n")
		ok = false
	}
	if Tweed.Password == "" {
		fmt.Fprintf(os.Stderr, "@R{(error)} No @R{--password} flag given, and @W{$TWEED_PASSWORD} not set.\n")
		ok = false
	}
	if !ok {
		os.Exit(1)
	}
}
