package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type UnbindCommand struct {
	NoWait bool `long:"no-wait" description:"don't wait for the binding to be created"`
	Args   struct {
		InstanceBinding []string `positional-arg-name:"instance/binding" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *UnbindCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	id, bid := GonnaNeedAnInstanceAndABinding(cmd.Args.InstanceBinding)

	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	var out api.UnbindResponse
	c.DELETE("/b/instances/"+id+"/bindings/"+bid, &out)

	if Tweed.JSON {
		JSON(out)
		if out.OK != "" {
			os.Exit(1)
		}
		os.Exit(0)
	}

	if out.OK != "" {
		if !Tweed.Quiet {
			fmt.Printf("@G{%s}\n", out.OK)
		}

		if !cmd.NoWait {
			await(c, patience{
				instance: id,
				task:     out.Ref,
				until:    "unbinding",
				negate:   true,
				quiet:    Tweed.Quiet,
			})
		}

		if !Tweed.Quiet {
			fmt.Printf("\nrun @W{tweed instance %s} for more details.\n", id)
		}
	} else {
		fmt.Printf("@R{%s}\n", out.Error)
		os.Exit(5)
	}
	return nil
}
