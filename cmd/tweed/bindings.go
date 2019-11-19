package main

import (
	"encoding/json"
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/jhunt/go-table"

	"github.com/tweedproject/tweed/api"
)

func Bindings(args []string) {
	GonnaNeedATweed()
	id := GonnaNeedAnInstance(args)

	var out api.BindingsResponse
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	c.GET("/b/instances/"+id+"/bindings", &out)

	if opts.JSON {
		JSON(out)
		os.Exit(0)
	}

	tbl := table.NewTable("Binding", "Credentials")
	for id, bi := range out.Bindings {
		b, err := json.MarshalIndent(bi, "", "  ")
		if err != nil {
			tbl.Row(nil, fmt.Sprintf("@C{%s}", id), fmt.Sprintf("@R{%s}", err))
		} else {
			tbl.Row(nil, fmt.Sprintf("@C{%s}", id), string(b))
		}
	}
	tbl.Output(os.Stdout)
}
