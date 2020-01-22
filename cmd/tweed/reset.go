package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"
)

type ResetCommand struct {
	Args struct {
		ID string `positional-arg-name:"instance" required:"true"`
	} `positional-args:"yes"`
	State string `short:"s" long:"state" default:"quiet" description:"expected state"`
}

func (cmd *ResetCommand) Execute(args []string) error {
	GonnaNeedATweed()

	switch cmd.State {
	case "quiet":
	case "provisioning":
	case "binding":
	case "unbinding":
	case "deprovisioning":
	case "gone":

	default:
		fmt.Fprintf(os.Stderr, "@R{(error)} invalid @Y{--state '%s'}\n")
		fmt.Fprintf(os.Stderr, "        only the following values are allowed:\n\n")
		fmt.Fprintf(os.Stderr, "          - @W{quiet}           nothing is happening\n")
		fmt.Fprintf(os.Stderr, "          - @W{provisioning}    new instance is being set up\n")
		fmt.Fprintf(os.Stderr, "          - @W{binding}         new credentials are being bound\n")
		fmt.Fprintf(os.Stderr, "          - @W{unbinding}       credentials are being unbound\n")
		fmt.Fprintf(os.Stderr, "          - @W{deprovisioning}  instance is being torn down\n")
		fmt.Fprintf(os.Stderr, "          - @W{gone}            instance deprovisioned\n")
		fmt.Fprintf(os.Stderr, "\n")
		os.Exit(1)
	}

	id := cmd.Args.ID
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	return c.PUT("/b/instances/"+id+"/state/"+cmd.State, nil, nil)
}
