package main

import (
	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

func Log(args []string) {
	GonnaNeedATweed()
	id := GonnaNeedAnInstance(args)

	var out api.InstanceResponse
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	c.GET("/b/instances/"+id, &out)
	fmt.Printf("%s\n", out.Log)
}
