package main

import (
	"os"

	"github.com/jhunt/go-table"

	"github.com/tweedproject/tweed/api"
)

func Instances(args []string) {
	GonnaNeedATweed()
	DontWantNoArgs(args)

	var out []api.InstanceResponse
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	err := c.GET("/b/instances", &out)
	bail(err)

	if opts.JSON {
		JSON(out)
		os.Exit(0)
	}

	tbl := table.NewTable("ID", "State", "Service", "Plan")
	for _, inst := range out {
		tbl.Row(nil, inst.ID, inst.State, inst.Service, inst.Plan)
	}
	tbl.Output(os.Stdout)
}
