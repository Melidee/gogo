package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Melidee/gogo/cli"
	"github.com/go-git/go-git/v6"
)

func main() {
	err := Command().Apply(os.Args)
	if err != nil {
		panic(err.Error())
	}
}

// resolve returns the value of a configuration flag from a variety of sources.
// The symbol to resolve should be given in kebab-case
// 
// The precedence of sources is as follows:
// Command line flags, environment variables, configuration values, then default values.
func resolve(symbol string) string {
	// env variables are typically in SCREAMING_SNAKE_CASE, so we convert the 
	// kebab-case symbol to SCREAMING_SNAKE_CASE
	screamingSnake := kebabToScreamingSnake(symbol)
	env := os.Getenv(screamingSnake)
	return env
}

func kebabToScreamingSnake(s string) string {
	upper := strings.ToUpper(s)
	return strings.Replace(upper, "-", "_", -1)
}

func Command() *cli.Command[cli.Empty] {
	return cli.NewCommand("gogo", cli.Empty{}).
		About("A simple CLI tool for Go.").
		Version("0.1.0").
		Author("Melidee <github.com/Melidee>").
		Usage("gogo [COMMAND] [OPTIONS]...").
		Action(func(ctx cli.Context[cli.Empty], value string) error { return nil }).
		Subcommand(SearchCommand()).
		Subcommand(InitCommand)
}

func SearchCommand() *cli.Command[Search] {
	return cli.NewCommand("search", NewSearch()).
		About("Search for packages in the go package repository.").
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

func (i Init) Init(pkgURL string) error {
	if pkgURL == "" {
		panic("no package name")
	}
	os.Mkdir(pkgURL, 0755)
	os.Chdir(pkgURL)

	if i.gitInit {
		git.PlainInit(".", false, git.WithDefaultBranch("main"))
	}
	makeGoMod(pkgURL)
	return nil
}

func makeGoMod(pkgName string) error {
	_, err := os.Stat("go.mod")
	if !errors.Is(err, os.ErrNotExist) {
		return errors.New("")
	}
	
	goVersion := runtime.Version()
	if goVersion == "" || goVersion == "unknown" {
		panic("unknown go version")
	}
	goVersion = goVersion[2:] // remove "go" prefix

	contents := fmt.Sprintf("module %s\n\ngo %s", pkgName, goVersion)
	f, err := os.Create("go.mod")
	if err != nil {
		panic("failed to create go mod file")
	}
	f.WriteString(contents)
	return nil
}

