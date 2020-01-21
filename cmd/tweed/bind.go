package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
	"github.com/tweedproject/tweed/random"
)

type BindCommand struct {
	ID     string `long:"as" optional:"yes" description:"use given binding id otherwise use random"`
	NoWait bool   `long:"no-wait" description:"don't wait for the binding to be created"`

	Args struct {
		ID string `positional-arg-name:"instance" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *BindCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()

	bid := cmd.ID
	if bid == "" {
		bid = random.ID("b")
	}

	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	var out api.BindResponse
	c.PUT("/b/instances/"+cmd.Args.ID+"/bindings/"+bid, nil, &out)

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
				instance: cmd.Args.ID,
				task:     out.Ref,
				until:    "binding",
				negate:   true,
				quiet:    Tweed.Quiet,
			})
		}

		if Tweed.Quiet {
			fmt.Printf("%s\n", bid)
		} else {
			fmt.Printf("binding: @C{%s}\n\n", bid)
			fmt.Printf("run @W{tweed instance %s} for more details.\n", cmd.Args.ID)
			fmt.Printf("run @W{tweed bindings %s} to show all bound credentials.\n", cmd.Args.ID)
		}
	} else {
		fmt.Printf("@R{%s}\n", out.Error)
		os.Exit(5)
	}
	return nil
}
