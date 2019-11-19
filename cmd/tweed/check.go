package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/tweedproject/tweed/api"
)

func Check(args []string) {
	GonnaNeedATweed()
	service, plan := GonnaNeedAServiceAndAPlan(args)

	var out api.ViabilityResponse
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	c.GET("/b/catalog/"+service+"/"+plan, &out)

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
	} else {
		fmt.Printf("@R{%s}\n", out.Error)
		os.Exit(5)
	}
}
