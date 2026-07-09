package cli

import "fmt"

func SearchCommand() *Command[Search] {
	return NewCommand("search", NewSearch()).
		SetHelp("Search for packages in the go package repository.").
		Action(func(ctx Context[Search], value string) {
			search := ctx.State()
			search.Query = value
			ctx.State().Search()
		}).
		AddFlag(NewFlag[Search]("count").
			SetShort('c').
			SetLong("count").
			SetAbout("Limit of search results to return").
			SetDefault("5").
			ActionSetInt(func(state Search) *int { return &state.Count })).
		AddFlag(NewFlag[Search]("filter").
			SetShort('f').
			SetLong("filter").
			SetAbout("Filter results by regular expression").
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

func Example() *Command[Empty] {
	return NewCommand("gogo", Empty{}).
		SetAuthor("Melidee").
		SetVersion("0.1.0").
		SetHelp("A simple CLI tool").
		AddFlag(NewFlag[Empty]("").Action(func(ctx Context[Empty], value string) {})).
		AddSubcommand(SearchCommand())
}
