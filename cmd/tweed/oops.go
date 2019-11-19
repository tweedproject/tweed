package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/tweedproject/tweed/api"
)

func Oops(args []string) {
	GonnaNeedATweed()
	id := GonnaNeedAnError(args)

	var out api.OopsResponse
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	c.GET("/b/oops/"+id, &out)

	if opts.JSON {
		JSON(out)
		os.Exit(0)
	}

	fmt.Printf("id:       %s\n", out.ID)
	fmt.Printf("client:   %s\n", out.Remote)
	fmt.Printf("handler:  %s\n", out.Handler)
	fmt.Printf("dated:    %s\n", out.Dated)
	fmt.Printf("message:  %s\n", out.Message)
	fmt.Printf("\n%s\n", out.Request)
}
