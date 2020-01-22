package main

import (
	"os"

	"github.com/jhunt/go-table"

	"github.com/tweedproject/tweed/api"
)

type FilesCommand struct {
	Args struct {
		ID string `positional-arg-name:"instance" required:"true"`
	} `positional-args:"yes"`
}

func (cmd *FilesCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	id := cmd.Args.ID

	var out []api.FileResponse
	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	c.GET("/b/instances/"+id+"/files", &out)

	if Tweed.JSON {
		JSON(out)
		os.Exit(0)
	}

	tbl := table.NewTable("File", "Summary", "Description")
	for _, file := range out {
		tbl.Row(nil, file.Filename, file.Summary, file.Description)
	}
	tbl.Output(os.Stdout)
	return nil
}
