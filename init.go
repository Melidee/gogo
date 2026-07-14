package main

import (
	"fmt"

	"github.com/Melidee/gogo/cli"
)

var NewCommand = cli.NewCommand("new", NewInit()).
	About("Initialize a new go project").
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
	About("Initialize a new go project").
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
