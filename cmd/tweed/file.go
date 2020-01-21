package main

import (
	"os"

	fmt "github.com/jhunt/go-ansi"

	"github.com/tweedproject/tweed/api"
)

type FileCommand struct {
	Args struct {
		ID   string `positional-arg-name:"instance" required:"true"`
		File string `positional-arg-name:"file-name" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *FileCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	id := cmd.Args.ID
	name := cmd.Args.File

	var out []api.FileResponse
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	c.GET("/b/instances/"+id+"/files", &out)

	if Tweed.JSON {
		JSON(out)
		os.Exit(0)
	}

	for _, file := range out {
		if file.Filename == name {
			fmt.Printf("%s", file.Contents)
			os.Exit(0)
		}
	}
	fmt.Fprintf(os.Stderr, "@R{%s}: file not found\n", name)
	return nil
}
