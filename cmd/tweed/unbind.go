package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/tweedproject/tweed/api"
)

func Unbind(args []string) {
	GonnaNeedATweed()
	id, bid := GonnaNeedAnInstanceAndABinding(args)

	c := Connect(opts.Tweed, opts.Username, opts.Password)
	var out api.UnbindResponse
	c.DELETE("/b/instances/"+id+"/bindings/"+bid, &out)

	if opts.JSON {
		JSON(out)
		if out.OK != "" {
			os.Exit(1)
		}
		os.Exit(0)
	}

	if out.OK != "" {
		if !opts.Quiet {
			fmt.Printf("@G{%s}\n", out.OK)
		}

		if opts.Unbind.Wait {
			await(c, patience{
				instance: id,
				until:    "unbinding",
				negate:   true,
				quiet:    opts.Quiet,
			})
		}

		if !opts.Quiet {
			fmt.Printf("\nrun @W{tweed instance %s} for more details.\n", id)
		}
	} else {
		fmt.Printf("@R{%s}\n", out.Error)
		os.Exit(5)
	}
}
