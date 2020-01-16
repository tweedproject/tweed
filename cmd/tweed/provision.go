package main

import (
	"os"
	"strings"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
	"github.com/tweedproject/tweed/random"
)

func Provision(args []string) {
	GonnaNeedATweed()
	service, plan := GonnaNeedAServiceAndAPlan(args)

	id := opts.Provision.ID
	if id == "" {
		id = random.ID("i")
	}

	params := make(map[string]interface{})
	for _, p := range opts.Provision.Params {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 1 {
			params[kv[0]] = ""
		} else {
			params[kv[0]] = kv[1]
		}
	}
	in := api.ProvisionRequest{
		Service: service,
		Plan:    plan,
		Params:  params,
	}

	c := Connect(opts.Tweed, opts.Username, opts.Password)
	var out api.ProvisionResponse
	c.PUT("/b/instances/"+id, in, &out)

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

		if opts.Provision.Wait {
			await(c, patience{
				instance: id,
				task:     out.Ref,
				until:    "provisioning",
				negate:   true,
				quiet:    opts.Quiet,
			})
		}

		if opts.Quiet {
			fmt.Printf("%s\n", id)
		} else {
			fmt.Printf("instance: @C{%s}\n\n", id)
			fmt.Printf("run @W{tweed instance %s} for more details.\n", id)
			fmt.Printf("run @W{tweed bind %s} to get some credentials.\n", id)
		}
	} else {
		fmt.Printf("@R{%s}\n", out.Error)
		os.Exit(5)
	}
}
