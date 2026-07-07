package hopeful

import "strconv"

type Search struct {
	Count int
	Filter string
}

func NewSearch() Search {
	return Search{
		Count: 5,
		Filter: "",
	}
}

type Add struct {
	
}

type Init struct {
	
}

func Test() {
	NewCmd("gogo").
		SetAuthor("Melidee").
		SetVersion("0.1.0").
		SetHelp("A simple CLI tool").
		AddFlag(NewFlag_("help").
			SetShort('h').
			SetLong("help").
			SetAbout("Show help message")).
		AddSubcommand(NewCommand("search", NewSearch()).
			SetHelp("Search for packages in the go package repository.").
			AddFlag(NewFlag[Search]("count").
				SetShort('c').
				SetLong("count").
				SetAbout("Limit of search results to return").
				SetDefault("5").
				Action(func(ctx Context[Search], value string) {
					count, err := strconv.Atoi(value)
					if err != nil {
						panic(err.Error())
					}
					ctx.State().Count = count
				})).
			AddFlag(NewFlag[Search]("filter").
				SetShort('f').
				SetLong("filter").
				SetAbout("Filter results by regular expression").
				Action(func(ctx Context[Search], value string) {
					ctx.State().Filter = value
				}))).
		AddSubcommand(NewCommand_("init").
			SetHelp("Initialize a new go project").
			AddFlag(NewFlag_("lib").SetAbout("")))
}