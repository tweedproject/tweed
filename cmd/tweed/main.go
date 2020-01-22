package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
)

var (
	Version     string
	BuildNumber string
	Tweed       TweedCommand
)

func version(prefix string) string {
	if Version == "" {
		return "(development)"
	} else {
		return prefix + Version
	}
}

func build() string {
	return BuildNumber
}

func main() {
	Tweed.Version = func() {
		fmt.Printf("tweed %s %s\n", version(""), build())
		os.Exit(0)
	}

	var parser = flags.NewParser(&Tweed, flags.Default)
	parser.NamespaceDelimiter = "-"

	Tweed.Broker.WireDynamicFlags(parser.Command.Find("broker"))

	_, err := parser.Parse()
	handleError(parser, err)

}

func handleError(helpParser *flags.Parser, err error) {
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}

		os.Exit(1)
	}
}
