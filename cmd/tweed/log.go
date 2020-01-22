package main

import (
	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type LogCommand struct {
	Args struct {
		ID string `positional-arg-name:"instance" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *LogCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	id := cmd.Args.ID

	var out api.InstanceResponse
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	c.GET("/b/instances/"+id, &out)
	fmt.Printf("%s\n", out.Log)
	return nil
}
