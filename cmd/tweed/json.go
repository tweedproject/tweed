package main

import (
	"encoding/json"
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func JSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{(error)} failed to marshal JSON:\n")
		fmt.Fprintf(os.Stderr, "        @R{%s}\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s\n", string(b))
}
