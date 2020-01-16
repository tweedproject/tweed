package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

func Deprovision(args []string) {
	GonnaNeedATweed()
	ids := GonnaNeedAtLeastOneInstance(args)
	c := Connect(opts.Tweed, opts.Username, opts.Password)

	rc := 0
	for _, id := range ids {
		if rc1 := deprovision1(c, id); rc1 > rc {
			rc = rc1
		}
	}
	os.Exit(rc)
}

func deprovision1(c *client, id string) int {
	var out api.DeprovisionResponse
	c.DELETE("/b/instances/"+id, &out)

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

		if opts.Deprovision.Wait {
			await(c, patience{
				instance: id,
				task:     out.Ref,
				until:    "gone",
				quiet:    opts.Quiet,
			})
		}

		if !opts.Quiet {
			fmt.Printf("\nrun @W{tweed instance %s} for more historical information.\n", id)
			fmt.Printf("run @W{tweed purge %s} to remove all trace of this service instance.\n", id)
		}
		return 0
	}

	if out.Gone {
		if !opts.Quiet {
			fmt.Printf("@Y{%s.}\n", out.Error)
			fmt.Printf("run @W{tweed purge %s} to remove all trace of this service instance.\n", id)
		}
		return 0
	}

	fmt.Printf("@R{%s}\n", out.Error)
	return 5
}
