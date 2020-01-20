package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"
)

type WaitCommand struct {
	Instance bool   `short:"i" long:"instance" description:"wait for an instance"`
	Task     bool   `short:"t" long:"task"     description:"wait for a task"`
	State    string `short:"s" long:"state" default:"quiet" description:"expected state"`
	Not      bool   `long:"not" description:"wait until out of given state"`
	Sleep    int    `short:"S" long:"sleep" default:"2" description:"number of seconds to sleep"`
	Max      int    `short:"m" long:"max" default:"150" description:"max number of tries for waiting on state"`
	Args     struct {
		InstanceOrTask []string `positional-arg-name:"instance or task" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *WaitCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()

	if !cmd.Task && !cmd.Instance {
		fmt.Fprintf(os.Stderr, "@R{(error)} either of @C{--instance} (@C{-i}) or @C{--task} (@C{-t}) are required!\n")
		os.Exit(1)
	}
	if cmd.Task && cmd.Instance {
		fmt.Fprintf(os.Stderr, "@R{(error)} only one of @C{--instance} (@C{-i}) or @C{--task} (@C{-t}) is allowed!\n")
		os.Exit(1)
	}
	if cmd.Task && cmd.State != "" {
		fmt.Fprintf(os.Stderr, "@Y({warning)} the @Y{--state '%s'} argument is ignored in @C{--task} mode...\n", cmd.State)
	}

	if cmd.Instance && cmd.State == "" {
		cmd.State = "quiet"
	}
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

	if cmd.Task {
		id, tid := GonnaNeedAnInstanceAndATask(cmd.Args.InstanceOrTask)
		c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
		ok := await(c, patience{
			instance: id,
			task:     tid,
			nap:      cmd.Sleep,
			max:      cmd.Max,
			quiet:    Tweed.Quiet,
		})
		if ok {
			os.Exit(0)
		} else {
			os.Exit(5)
		}

	} else {
		id := GonnaNeedAnInstance(cmd.Args.InstanceOrTask)
		c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
		ok := await(c, patience{
			instance: id,
			until:    cmd.State,
			negate:   cmd.Not,
			nap:      cmd.Sleep,
			max:      cmd.Max,
			quiet:    Tweed.Quiet,
		})
		if ok {
			os.Exit(0)
		} else {
			os.Exit(5)
		}
	}
	return nil
}
