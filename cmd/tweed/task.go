package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type TaskCommand struct {
	NoWait bool `long:"no-wait" description:"don't wait for the task to be 'quiet'"`
	Args   struct {
		InstanceTask []string `positional-arg-name:"instance/task" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *TaskCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	id, tid := GonnaNeedAnInstanceAndATask(cmd.Args.InstanceTask)

	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	if !cmd.NoWait {
		await(c, patience{
			instance: id,
			task:     tid,
			quiet:    Tweed.Quiet,
		})
	}

	var out api.TaskResponse
	err := c.GET("/b/instances/"+id+"/tasks/"+tid, &out)

	if err != nil {
		fmt.Printf("failed: %s", err)
		os.Exit(1)
	}

	if Tweed.JSON {
		JSON(out)
		os.Exit(0)
	}

	fmt.Printf("task: %s\n", out.Task)
	if out.Done {
		if out.Exited {
			fmt.Printf("stat: exited %d\n", out.ExitCode)
		} else {
			fmt.Printf("stat: terminated\n")
		}
	} else {
		fmt.Printf("stat: running\n")
	}
	fmt.Printf("---[ stdout ]-------------\n%s\n\n", out.Stdout)
	fmt.Printf("---[ stderr ]-------------\n%s\n\n", out.Stderr)
	fmt.Printf("--------------------------\n\n")
	return nil
}
