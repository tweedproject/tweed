package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type CheckCommand struct {
	Args struct {
		ServicePlan []string `positional-arg-name:"service/plan" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *CheckCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	service, plan := GonnaNeedAServiceAndAPlan(cmd.Args.ServicePlan)

	var out api.ViabilityResponse
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	c.GET("/b/catalog/"+service+"/"+plan, &out)

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
	} else {
		fmt.Printf("@R{%s}\n", out.Error)
		os.Exit(5)
	}
	return nil
}
