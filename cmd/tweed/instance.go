package main

import (
	"encoding/json"
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/tweedproject/tweed/api"
)

func Instance(args []string) {
	GonnaNeedATweed()
	id := GonnaNeedAnInstance(args)

	var out api.InstanceResponse
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	c.GET("/b/instances/"+id, &out)

	if opts.JSON {
		JSON(out)
		os.Exit(0)
	}

	fmt.Printf("id:       %s\n", out.ID)
	fmt.Printf("state:    %s\n", out.State)
	fmt.Printf("service:  %s\n", out.Service)
	fmt.Printf("plan:     %s\n", out.Plan)
	if out.Params != nil {
		fmt.Printf("params:\n")
		for k, v := range out.Params {
			fmt.Printf("  %s = %v\n", k, v)
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("params:   (none)\n")
	}

	if out.Bindings == nil {
		fmt.Printf("bindings: (none)\n")
	} else {
		fmt.Printf("bindings:\n")
		for id, bi := range out.Bindings {
			b, err := json.MarshalIndent(bi, "    ", "  ")
			if err != nil {
				fmt.Printf("  %s: @R{%s}\n", id, err)
			} else {
				fmt.Printf("  %s:\n    %s", id, string(b))
			}
			fmt.Printf("\n")
		}
	}
}
