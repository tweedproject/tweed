package main

import (
	"os"

	"github.com/jhunt/go-table"

	"github.com/tweedproject/tweed/api"
)

func Files(args []string) {
	GonnaNeedATweed()
	id := GonnaNeedAnInstance(args)

	var out []api.FileResponse
	c := Connect(opts.Tweed, opts.Username, opts.Password)
	c.GET("/b/instances/"+id+"/files", &out)

	if opts.JSON {
		JSON(out)
		os.Exit(0)
	}

	tbl := table.NewTable("File", "Summary", "Description")
	for _, file := range out {
		tbl.Row(nil, file.Filename, file.Summary, file.Description)
	}
	tbl.Output(os.Stdout)
}
