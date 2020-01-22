package main

import (
	"encoding/json"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type BindingCommand struct {
	Args struct {
		Binding []string `positional-arg-name:"instance/binding" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *BindingCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	id, bid := GonnaNeedAnInstanceAndABinding(cmd.Args.Binding)

	var out api.BindingResponse
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	c.GET("/b/instances/"+id+"/bindings/"+bid, &out)

	b, err := json.MarshalIndent(out.Binding, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(b))
	return nil
}
