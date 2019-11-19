package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"
)

func Wait(args []string) {
	GonnaNeedATweed()

	if !opts.Wait.Task && !opts.Wait.Instance {
		fmt.Fprintf(os.Stderr, "@R{(error)} either of @C{--instance} (@C{-i}) or @C{--task} (@C{-t}) are required!\n")
		os.Exit(1)
	}
	if opts.Wait.Task && opts.Wait.Instance {
		fmt.Fprintf(os.Stderr, "@R{(error)} only one of @C{--instance} (@C{-i}) or @C{--task} (@C{-t}) is allowed!\n")
		os.Exit(1)
	}
	if opts.Wait.Task && opts.Wait.State != "" {
		fmt.Fprintf(os.Stderr, "@Y({warning)} the @Y{--state '%s'} argument is ignored in @C{--task} mode...\n", opts.Wait.State)
	}

	if opts.Wait.Instance && opts.Wait.State == "" {
		opts.Wait.State = "quiet"
	}
	switch opts.Wait.State {
	case "quiet":
	case "provisioning":
	case "binding":
	case "unbinding":
	case "deprovisioning":
	case "gone":

	default:
		fmt.Fprintf(os.Stderr, "@R{(error)} invalid @Y{--state '%s'}\n")
		fmt.Fprintf(os.Stderr, "        only the following values are allowed:\n\n")
		fmt.Fprintf(os.Stderr, "          - @W{quiet}           nothing is happening\n")
		fmt.Fprintf(os.Stderr, "          - @W{provisioning}    new instance is being set up\n")
		fmt.Fprintf(os.Stderr, "          - @W{binding}         new credentials are being bound\n")
		fmt.Fprintf(os.Stderr, "          - @W{unbinding}       credentials are being unbound\n")
		fmt.Fprintf(os.Stderr, "          - @W{deprovisioning}  instance is being torn down\n")
		fmt.Fprintf(os.Stderr, "          - @W{gone}            instance deprovisioned\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	if opts.Wait.Task {
		id, tid := GonnaNeedAnInstanceAndATask(args)
		c := Connect(opts.Tweed, opts.Username, opts.Password)
		ok := await(c, patience{
			instance: id,
			task:     tid,
			nap:      opts.Wait.Sleep,
			max:      opts.Wait.Max,
			quiet:    opts.Quiet,
		})
		if ok {
			os.Exit(0)
		} else {
			os.Exit(5)
		}

	} else {
		id := GonnaNeedAnInstance(args)
		c := Connect(opts.Tweed, opts.Username, opts.Password)
		ok := await(c, patience{
			instance: id,
			until:    opts.Wait.State,
			negate:   opts.Wait.Not,
			nap:      opts.Wait.Sleep,
			max:      opts.Wait.Max,
			quiet:    opts.Quiet,
		})
		if ok {
			os.Exit(0)
		} else {
			os.Exit(5)
		}
	}
}
