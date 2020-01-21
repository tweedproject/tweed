package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type PurgeCommand struct {
	All  bool `short:"a" long:"all" description:"purge all instances which are in state 'gone'"`
	Args struct {
		InstanceIds []string `positional-arg-name:"instance"`
	} `positional-args:"yes"`
}

func (cmd *PurgeCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()

	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	ids := make([]string, 0)

	if cmd.All {
		var ls []api.InstanceResponse
		c.GET("/b/instances", &ls)

		for _, inst := range ls {
			if inst.State == "gone" {
				ids = append(ids, inst.ID)
			}
		}

	} else {
		ids = GonnaNeedAtLeastOneInstance(cmd.Args.InstanceIds)
	}

	rc := 0
	for _, id := range ids {
		if rc1 := purge1(c, id); rc1 > rc {
			rc = rc1
		}
	}
	os.Exit(rc)
	return nil
}

func purge1(c *client, id string) int {
	var out api.PurgeResponse
	c.DELETE("/b/instances/"+id+"/log", &out)

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
		return 0
	}
	fmt.Printf("@R{%s}\n", out.Error)
	return 5
}
