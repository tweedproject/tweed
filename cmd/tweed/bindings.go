package main

import (
	"encoding/json"
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/jhunt/go-table"

	"github.com/tweedproject/tweed/api"
)

type BindingsCommand struct {
	Args struct {
		ID string `positional-arg-name:"instance" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *BindingsCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	id := cmd.Args.ID

	var out api.BindingsResponse
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	c.GET("/b/instances/"+id+"/bindings", &out)

	if Tweed.JSON {
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
	return nil
}
