package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Melidee/gogo/cli"
)

func main() {
	err := Command().Apply(os.Args)
	if err != nil {
		panic(err.Error())
	}
}

func Command() *cli.Command[cli.Empty] {
	return cli.NewCommand("gogo", cli.Empty{}).
		About("A simple CLI tool for Go.").
		Version("0.1.0").
		Author("Melidee <github.com/Melidee>").
		Help("").
		Usage("gogo [COMMAND] [OPTIONS]...").
		Action(func(ctx cli.Context[cli.Empty], value string) error { return nil}).
		Subcommand(SearchCommand()).
		Subcommand(InitCommand())
}

func SearchCommand() *cli.Command[Search] {
	return cli.NewCommand("search", NewSearch()).
		About("Search for packages in the go package repository.").
		Help("").
		Action(func(ctx cli.Context[Search], query string) error {
			search := ctx.State()
			search.Query = query
			ctx.State().Search()
			return nil
		}).
		Flag(cli.NewFlag[Search]("count").
			Short('c').
			Long("count").
			About("Limit of search results to return").
			Default("5").
			ArgName("COUNT").
			ActionSetInt(func(state Search) *int { return &state.Count })).
		Flag(cli.NewFlag[Search]("filter").
			Short('f').
			Long("filter").
			About("Filter results by regular expression").
			ArgName("REGEX").
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
		gitInit:     true,
		packageName: "",
	}
}

func InitCommand() *cli.Command[Init] {
	return cli.NewCommand("init", NewInit()).
		About("Initialize a new go project").
		Help("").
		Action(func(ctx cli.Context[Init], pkgName string) error {
			if pkgName == "" {
				panic("no package name")
			}
			os.Mkdir(pkgName, 0755)
			os.Chdir("pkgName")
			if ctx.State().gitInit {
				exec.Command("git", "init", "-b", "main")
			}
			return nil
		}).
		Flag(cli.NewFlag[Init]("lib").
			Long("lib").
			About("Scaffold this project as a library, without an executable").
			ActionSetTrue(func(state Init) *bool { return &state.isLib })).
		Flag(cli.NewFlag[Init]("no-git").
			Long("no-git").
			About("Do not initialize a git repository").
			ActionSetFalse(func(state Init) *bool { return &state.gitInit }))

}
