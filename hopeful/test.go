package hopeful

type Search struct {
	Count  int
	Filter string
}

func NewSearch() Search {
	return Search{
		Count:  5,
		Filter: "",
	}
}

type Add struct {
}

func NewAdd() Add {
	return Add{}
}

type Init struct {
	isLib bool
}

func NewInit() Init {
	return Init{
		isLib: false,
	}
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
		AddSubcommand(NewCommand("init", NewInit()).
			SetHelp("Initialize a new go project").
			AddFlag(NewFlag[Init]("lib").
				SetLong("lib").
				SetAbout("initialize the project as a library").
				ActionToggleBool(func(state Init) *bool { return &state.isLib }))).
		AddSubcommand(NewCommand("search", NewSearch()).
			SetHelp("Search for packages in the go package repository.").
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
				ActionSet(func(state Search) *string { return &state.Filter }))).
		AddSubcommand(NewCommand("add", NewAdd()))
		
}
