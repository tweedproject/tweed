package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/tweedproject/tweed/api"
)

func Purge(args []string) {
	GonnaNeedATweed()
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	ids := make([]string, 0)

	if opts.Purge.All {
		var ls []api.InstanceResponse
		c.GET("/b/instances", &ls)

		for _, inst := range ls {
			if inst.State == "gone" {
				ids = append(ids, inst.ID)
			}
		}

	} else {
		ids = GonnaNeedAtLeastOneInstance(args)
	}

	rc := 0
	for _, id := range ids {
		if rc1 := purge1(c, id); rc1 > rc {
			rc = rc1
		}
	}
	os.Exit(rc)
}

func purge1(c *client, id string) int {
	var out api.PurgeResponse
	c.DELETE("/b/instances/"+id+"/log", &out)

	if opts.JSON {
		JSON(out)
		if out.OK != "" {
			return 1
		}
		return 0
	}

	if out.OK != "" {
		if !opts.Quiet {
			fmt.Printf("@G{%s}\n", out.OK)
		}
		return 0
	}
	fmt.Printf("@R{%s}\n", out.Error)
	return 5
}
