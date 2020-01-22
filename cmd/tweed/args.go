package main

import (
	"os"
	"strings"

	fmt "github.com/jhunt/go-ansi"
)

func DontWantNoArgs(args []string) {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "@R{(error)} too many arguments!\n")
		fmt.Fprintf(os.Stderr, "        (no clue what to do with @Y{%s})\n", strings.Join(args, " "))
		os.Exit(1)
	}
}

func gonnaNeedOneThing(what string, args []string) string {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "@R{(error)} missing @C{%s} argument!\n", what)
		os.Exit(1)

	} else if len(args) == 1 {
		return args[0]

	} else {
		fmt.Fprintf(os.Stderr, "@R{(error)} too many arguments!\n")
		fmt.Fprintf(os.Stderr, "        (no clue what to do with @Y{%s})\n", strings.Join(args[1:], " "))
		os.Exit(1)
	}

	return ""
}

func GonnaNeedAnInstance(args []string) string {
	return gonnaNeedOneThing("instance", args)
}

func gonnaNeedSomeThings(what string, args []string) []string {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "@R{(error)} missing @C{%s} argument!\n", what)
		os.Exit(1)
	}
	return args
}

func GonnaNeedAtLeastOneInstance(args []string) []string {
	return gonnaNeedSomeThings("instance", args)
}

func GonnaNeedAnError(args []string) string {
	return gonnaNeedOneThing("error", args)
}

func gonnaNeedTwoThings(first, second string, args []string) (string, string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "@R{(error)} missing @C{%s} and @C{%s} argument!\n", first, second)
		os.Exit(1)

	} else if len(args) == 1 {
		p := strings.Split(args[0], "/")
		if len(p) < 2 {
			fmt.Fprintf(os.Stderr, "@R{(error)} missing @C{%s} in @C{%s}@W{/}@C{%s} argument '@Y{%s}'!\n", second, first, second, args[0])
			os.Exit(1)
		}
		if len(p) > 2 {
			fmt.Fprintf(os.Stderr, "@R{(error)} too many slashes in @C{%s}@W{/}@C{%s} argument '@Y{%s}'!\n", first, second, args[0])
			os.Exit(1)
		}
		return strings.TrimSpace(p[0]), strings.TrimSpace(p[1])

	} else if len(args) == 2 {
		return args[0], args[1]

	} else {
		fmt.Fprintf(os.Stderr, "@R{(error)} too many arguments!\n")
		fmt.Fprintf(os.Stderr, "        (no clue what to do with @Y{%s})\n", strings.Join(args[2:], " "))
		os.Exit(1)
	}

	return "", ""
}

func GonnaNeedAServiceAndAPlan(args []string) (string, string) {
	return gonnaNeedTwoThings("service", "plan", args)
}

func GonnaNeedAnInstanceAndATask(args []string) (string, string) {
	return gonnaNeedTwoThings("instance", "task", args)
}

func GonnaNeedAnInstanceAndABinding(args []string) (string, string) {
	return gonnaNeedTwoThings("instance", "binding", args)
}

func GonnaNeedAnInstanceAndAFile(args []string) (string, string) {
	return gonnaNeedTwoThings("instance", "file", args)
}
