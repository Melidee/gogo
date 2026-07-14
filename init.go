package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/Melidee/gogo/cli"
	"github.com/go-git/go-git/v6"
)

var NewCommand = cli.NewCommand("new", NewInit()).
	About("Create a new go project").
	Action(func(ctx cli.Context[Init], pkgName string) error {
		fmt.Printf("Initializing new go project `%s`...\n", pkgName)
		return ctx.State().Init(pkgName)
	}).
	Flag(cli.NewFlag[Init]("lib").
		Long("lib").
		About("Scaffold this project as a library, without an executable").
		ActionSetTrue(func(state Init) *bool { return &state.isLib })).
	Flag(cli.NewFlag[Init]("no-git").
		Long("no-git").
		About("Do not initialize a git repository").
		ActionSetFalse(func(state Init) *bool { return &state.gitInit }))

var InitCommand = cli.NewCommand("init", NewInit()).
	About("Initialize a new go project in this directory").
	Action(func(ctx cli.Context[Init], pkgName string) error {
		fmt.Printf("Initializing new go project `%s`...\n", pkgName)
		return ctx.State().Init(pkgName)
	}).
	Flag(cli.NewFlag[Init]("lib").
		Long("lib").
		About("Scaffold this project as a library, without an executable").
		ActionSetTrue(func(state Init) *bool { return &state.isLib })).
	Flag(cli.NewFlag[Init]("no-git").
		Long("no-git").
		About("Do not initialize a git repository").
		ActionSetFalse(func(state Init) *bool { return &state.gitInit }))

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
