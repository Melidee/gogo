package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	cmd := gogoCmd()
	matches := cmd.Apply(os.Args)
	if matches.HasFlag("help") {
		fmt.Print(cmd.HelpMessage(""))
	}
	if matches.HasFlag("version") {
		fmt.Printf("gogo v%s", cmd.Version)
	}
	if matches.Subcommand == nil {
		return
	}
	switch matches.Subcommand.Name {
	case "search":
		log.Printf("search command matched with flags: %+v\n", matches.Flags)
	default:
		log.Printf("no command matched, args: %+v\n", os.Args)
	}
}

func gogoCmd[T any]() *Command[T] {
	return NewCommand("gogo").
		SetAuthor("Melidee").
		SetVersion("0.1.0").
		SetHelp("A simple CLI tool").
		AddFlag(NewFlag("help").
			SetShort('h').
			SetLong("help").
			SetAbout("Show help message")).
		AddSubcommand(NewCommand[any]("search").
			AddFlag(NewFlag("limit").
				SetShort('l').
				SetLong("count").
				SetAbout("Limit of search results to return").
				SetDefault("5")).
			AddFlag(NewFlag("filter").
				SetShort('f').
				SetLong("filter").
				SetAbout("Filter results by regular expression")).
		AddSubcommand(NewStatelessCommand("init").
			AddFlag(NewFlag("lib").SetLong("lib"))))
}



func NewStatelessCommand(name string) *Command[struct{}] {
	return &Command[struct{}]{
		Name: name,
	}
}

func NewCommand[T any](name string, initialState T) *Command[T] {
	return &Command[T]{
		Name: name,
		State: initialState,
	}
}

func (c *Command[T]) SetAuthor(author string) *Command[T] {
	c.Author = author
	return c
}

func (c *Command[T]) SetVersion(version string) *Command[T] {
	versionFlag := NewFlag("version").
		SetShort('v').
		SetLong("version").
		SetAbout("Show program version").
		SetAction(func(_ string) {
			fmt.Printf("%s v%s", c.Name, version)
			os.Exit(0)
		})

	c.Flags = append(c.Flags, versionFlag)
	c.Version = version
	return c
}

func (c *Command[T]) SetHelp(help string) *Command[T] {
	c.Help = help
	return c
}

func (c *Command[T]) AddFlag(flag *Flag) *Command[T] {
	c.Flags = append(c.Flags, flag)
	return c
}

func (c *Command[T]) AddSubcommand(cmd *Command[any]) *Command[T] {
	c.Subcommands = append(c.Subcommands, cmd)
	return c
}

func (c *Command) HelpMessage(_ string) string {
	return fmt.Sprintf("%s v%s\n%s\n%s\n", c.Name, c.Version, c.Author, c.About)
}

func (c *Command) Apply(args []string) *CmdMatches {
	iter := NewArgIterator(args)
	return c.apply(iter)
}

func (c *Command) apply(iter *ArgIterator) *CmdMatches {
	matches := NewCmdMatches(iter.Next())

	if iter.NextIsFlag() {
		for _, flag := range c.Flags {
			if match := flag.Matches(iter); match != nil {
				matches.Flags = append(matches.Flags, *match)
				break
			}
		}
	} else {
		// next is subcommand
		for _, sub := range c.Subcommands {
			if sub.Name == iter.Peek() {
				matches.Subcommand = sub.apply(iter)
				break
			}
		}
	}

	return matches
}

type CmdMatches struct {
	Name       string
	Flags      []FlagMatch
	Subcommand *CmdMatches
}

func NewCmdMatches(name string) *CmdMatches {
	return &CmdMatches{
		Name: name,
	}
}

func (c *CmdMatches) HasFlag(name string) bool {
	for _, flag := range c.Flags {
		if flag.Name == name {
			return true
		}
	}
	return false
}

type ArgIterator struct {
	cursor int
	args   []string
}

func NewArgIterator(args []string) *ArgIterator {
	return &ArgIterator{
		cursor: 0,
		args:   args,
	}
}

func (a *ArgIterator) HasNext() bool {
	return a.cursor < len(a.args)
}

func (a *ArgIterator) Next() string {
	if a.cursor >= len(a.args) {
		return ""
	}
	a.cursor++
	return a.args[a.cursor-1]
}

func (a *ArgIterator) Peek() string {
	if a.cursor >= len(a.args) {
		return ""
	}
	return a.args[a.cursor]
}

func (a *ArgIterator) NextIsFlag() bool {
	if !a.HasNext() {
		return false
	}
	arg := a.Peek()
	return len(arg) > 0 && arg[0] == '-'
}

func (a *ArgIterator) NextIsSubcommand() bool {
	return a.HasNext() && !a.NextIsFlag()
}
