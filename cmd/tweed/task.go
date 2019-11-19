package main

import (
	fmt "github.com/jhunt/go-ansi"
	"os"

	"github.com/tweedproject/tweed/api"
)

func Task(args []string) {
	GonnaNeedATweed()
	id, tid := GonnaNeedAnInstanceAndATask(args)

	c := Connect(opts.Tweed, opts.Username, opts.Password)
	if opts.Task.Wait {
		await(c, patience{
			instance: id,
			task:     tid,
			quiet:    opts.Quiet,
		})
	}

	var out api.TaskResponse
	c.GET("/b/instances/"+id+"/tasks/"+tid, &out)

	if opts.JSON {
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
}
