package main

import (
	"os"

	"github.com/jhunt/go-table"

	"github.com/tweedproject/tweed/api"
)

type InstancesCommand struct {
}

func (cmd *InstancesCommand) Execute(args []string) error {
	GonnaNeedATweed()
	DontWantNoArgs(args)

	var out []api.InstanceResponse
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	err := c.GET("/b/instances", &out)
	bail(err)

	if Tweed.JSON {
		JSON(out)
		os.Exit(0)
	}

	tbl := table.NewTable("ID", "State", "Service", "Plan")
	for _, inst := range out {
		tbl.Row(nil, inst.ID, inst.State, inst.Service, inst.Plan)
	}
	tbl.Output(os.Stdout)
	return nil
}
