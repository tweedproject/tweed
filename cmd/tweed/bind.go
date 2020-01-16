package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
	"github.com/tweedproject/tweed/random"
)

func Bind(args []string) {
	GonnaNeedATweed()
	id := GonnaNeedAnInstance(args)

	bid := opts.Bind.ID
	if bid == "" {
		bid = random.ID("b")
	}

	c := Connect(opts.Tweed, opts.Username, opts.Password)
	var out api.BindResponse
	c.PUT("/b/instances/"+id+"/bindings/"+bid, nil, &out)

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

		if opts.Bind.Wait {
			await(c, patience{
				instance: id,
				task:     out.Ref,
				until:    "binding",
				negate:   true,
				quiet:    opts.Quiet,
			})
		}

		if opts.Quiet {
			fmt.Printf("%s\n", bid)
		} else {
			fmt.Printf("binding: @C{%s}\n\n", bid)
			fmt.Printf("run @W{tweed instance %s} for more details.\n", id)
			fmt.Printf("run @W{tweed bindings %s} to show all bound credentials.\n", id)
		}
	} else {
		fmt.Printf("@R{%s}\n", out.Error)
		os.Exit(5)
	}
}
