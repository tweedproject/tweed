package main

import (
	"os"
	"strings"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
	"github.com/tweedproject/tweed/random"
)

type ProvisionCommand struct {
	ID     string   `long:"as" optional:"yes" description:"use given service id otherwise use random"`
	NoWait bool     `long:"no-wait" description:"don't wait for the service to be created"`
	Params []string `short:"P" optional:"yes" long:"params" description:"params passed to the service"`
	Args   struct {
		ServicePlan []string `positional-arg-name:"service/plan" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *ProvisionCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	service, plan := GonnaNeedAServiceAndAPlan(cmd.Args.ServicePlan)

	id := cmd.ID
	if id == "" {
		id = random.ID("i")
	}

	params := make(map[string]interface{})
	for _, p := range cmd.Params {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 1 {
			params[kv[0]] = ""
		} else {
			params[kv[0]] = kv[1]
		}
	}
	in := api.ProvisionRequest{
		Service: service,
		Plan:    plan,
		Params:  params,
	}

	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	var out api.ProvisionResponse
	c.PUT("/b/instances/"+id, in, &out)

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
				until:    "provisioning",
				negate:   true,
				quiet:    Tweed.Quiet,
			})
		}

		if !Tweed.Quiet {
			fmt.Printf("%s\n", id)
		} else {
			fmt.Printf("instance: @C{%s}\n\n", id)
			fmt.Printf("run @W{tweed instance %s} for more details.\n", id)
			fmt.Printf("run @W{tweed bind %s} to get some credentials.\n", id)
		}
	} else {
		fmt.Printf("@R{%s}\n", out.Error)
		os.Exit(5)
	}
	return nil
}
