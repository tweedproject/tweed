package main

import (
	"os"
	"strings"

	fmt "github.com/jhunt/go-ansi"

	"github.com/jhunt/go-table"

	"github.com/tweedproject/tweed"
)

type CatalogCommand struct {
}

func (cmd *CatalogCommand) Execute(args []string) error {
	SetupLogging()
	GonnaNeedATweed()
	DontWantNoArgs(args)

	c := Connect(Tweed.Tweed, Tweed.Username, Tweed.Password)
	var cat tweed.Catalog
	c.GET("/b/catalog", &cat)

	if Tweed.JSON {
		JSON(cat)
		os.Exit(0)
	}

	tbl := table.NewTable("Service / Plan", "#", "Free?", "")
	for _, s := range cat.Services {
		tags := "(none)"
		if len(s.Tags) > 0 {
			t := make([]string, len(s.Tags))
			for i := range s.Tags {
				t[i] = fmt.Sprintf("@C{%s}", s.Tags[i])
			}
			tags = strings.Join(t, ", ")
		}

		for _, p := range s.Plans {
			free := "no"
			if p.Free {
				free = fmt.Sprintf("@G{yes}")
			}

			tbl.Row(nil,
				fmt.Sprintf("@W{%s}/@W{%s}", s.ID, p.ID),
				fmt.Sprintf("%d/%d", p.Tweed.Provisioned, p.Tweed.Limit),
				free,
				fmt.Sprintf("%s\n@W{%s}\n[tags: %s]\n\n",
					strings.TrimSpace(s.Description),
					strings.TrimSpace(p.Description),
					tags),
			)
		}
	}
	tbl.Output(os.Stdout)
	return nil
}
