package main

import (
	"fmt"

	"github.com/Melidee/gogo/cli"
)

func Command() *cli.Command[cli.Empty] {
	return cli.NewCommand("gogo", cli.Empty{}).
		Help("A simple CLI tool for Go.").
		Version("0.1.0").
		Author("Melidee").
		Subcommand(SearchCommand())
}

func SearchCommand() *cli.Command[Search] {
	return cli.NewCommand("search", NewSearch()).
		Help("Search for packages in the go package repository.").
		Action(func(ctx cli.Context[Search], value string) {
			search := ctx.State()
			search.Query = value
			ctx.State().Search()
		}).
		Flag(cli.NewFlag[Search]("count").
			Short('c').
			Long("count").
			About("Limit of search results to return").
			Default("5").
			ActionSetInt(func(state Search) *int { return &state.Count })).
		Flag(cli.NewFlag[Search]("filter").
			Short('f').
			Long("filter").
			About("Filter results by regular expression").
			ActionSet(func(state Search) *string { return &state.Filter }))
}

type Search struct {
	Count  int
	Filter string
	Query  string
}

func NewSearch() Search {
	return Search{
		Count:  5,
		Filter: "",
	}
}

func (s Search) Search() {
	fmt.Printf("searching `%s` with filter `%s` for %d results\n", s.Query, s.Filter, s.Count)
}

type Add struct {
	index int
}

func NewAdd() Add {
	return Add{}
}

type Init struct {
	isLib       bool
	gitInit     bool
	packageName string
}

func NewInit() Init {
	return Init{
		isLib:       false,
		gitInit:     false,
		packageName: "",
	}
}
