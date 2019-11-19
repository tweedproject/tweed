package main

import (
	"encoding/json"
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/tweedproject/tweed/api"
)

func Binding(args []string) {
	GonnaNeedATweed()
	id, bid := GonnaNeedAnInstanceAndABinding(args)

	var out api.BindingResponse
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	c.GET("/b/instances/"+id+"/bindings/"+bid, &out)

	b, err := json.MarshalIndent(out.Binding, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "@R{%s}", err)
		os.Exit(2)
	}
	fmt.Printf("%s\n", string(b))
}
