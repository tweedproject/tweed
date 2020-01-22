package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type OopsCommand struct {
	Args struct {
		ID string `positional-arg-name:"id" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *OopsCommand) Execute(args []string) error {
	GonnaNeedATweed()

	var out api.OopsResponse
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	c.GET("/b/oops/"+cmd.Args.ID, &out)

	if Tweed.JSON {
		JSON(out)
		os.Exit(0)
	}

	fmt.Printf("id:       %s\n", out.ID)
	fmt.Printf("client:   %s\n", out.Remote)
	fmt.Printf("handler:  %s\n", out.Handler)
	fmt.Printf("dated:    %s\n", out.Dated)
	fmt.Printf("message:  %s\n", out.Message)
	fmt.Printf("\n%s\n", out.Request)
	return nil
}
