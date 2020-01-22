package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type DeprovisionCommand struct {
	NoWait bool `long:"no-wait" description:"don't wait for the service be deprovisioned"`
	Args   struct {
		InstanceIds []string `positional-arg-name:"instance" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *DeprovisionCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	ids := GonnaNeedAtLeastOneInstance(cmd.Args.InstanceIds)
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)

	rc := 0
	for _, id := range ids {
		if rc1 := cmd.deprovision1(c, id); rc1 > rc {
			rc = rc1
		}
	}
	os.Exit(rc)
	return nil
}

func (cmd *DeprovisionCommand) deprovision1(c *client, id string) int {
	var out api.DeprovisionResponse
	c.DELETE("/b/instances/"+id, &out)

	if Tweed.JSON {
		JSON(out)
		if out.OK != "" {
			return 1
		}
		return 0
	}

	if out.OK != "" {
		if !Tweed.Quiet {
			fmt.Printf("@G{%s}\n", out.OK)
		}

		if !cmd.NoWait {
			await(c, patience{
				instance: id,
				task:     out.Ref,
				until:    "gone",
				quiet:    Tweed.Quiet,
			})
		}

		if !Tweed.Quiet {
			fmt.Printf("\nrun @W{tweed instance %s} for more historical information.\n", id)
			fmt.Printf("run @W{tweed purge %s} to remove all trace of this service instance.\n", id)
		}
		return 0
	}

	if out.Gone {
		if !Tweed.Quiet {
			fmt.Printf("@Y{%s.}\n", out.Error)
			fmt.Printf("run @W{tweed purge %s} to remove all trace of this service instance.\n", id)
		}
		return 0
	}

	fmt.Printf("@R{%s}\n", out.Error)
	return 5
}
