package main

import (
	"log"
	"os"
)

func main() {
	matches := gogoCmd().Matches(os.Args)
	if matches.HasFlag("help") {
		
	}
	if matches.HasFlag("version") {
		
	}
	switch matches.Subcommand.Name {
	case "search":
		log.Printf("search command matched with flags: %+v\n", matches.Flags)
	default:
		log.Printf("no command matched, args: %+v\n", os.Args)
	}
}

func gogoCmd() *Command {
	return NewCommand("gogo").
		SetAuthor("Melidee").
		SetVersion("0.1.0").
		SetHelp("A simple CLI tool").
		AddFlag(NewFlag("help").
			SetShort('h').
			SetLong("help").
			SetAbout("Show help message")).
		AddFlag(NewFlag("version").
			SetShort('v').
			SetLong("version").
			SetAbout("Show program version")).
		AddSubcommand(NewCommand("search").
			AddFlag(NewFlag("limit").
				SetShort('l').
				SetLong("count").
				SetAbout("Limit of search results to return").
				SetDefault("5")).
			AddFlag(NewFlag("filter").
				SetShort('f').
				SetLong("filter").
				SetAbout("Filter results by regular expression")))
}

type Command struct {
	Name        string
	Author      string
	Version     string
	Help        string
	Flags       []*Flag
	Subcommands []*Command
}

func NewCommand(name string) *Command {
	return &Command{
		Name: name,
	}
}

func (c *Command) SetAuthor(author string) *Command {
	c.Author = author
	return c
}

func (c *Command) SetVersion(version string) *Command {
	c.Version = version
	return c
}

func (c *Command) SetHelp(help string) *Command {
	c.Help = help
	return c
}

func (c *Command) AddFlag(flag *Flag) *Command {
	c.Flags = append(c.Flags, flag)
	return c
}

func (c *Command) AddSubcommand(cmd *Command) *Command {
	c.Subcommands = append(c.Subcommands, cmd)
	return c
}

func (c *Command) Matches(args []string) *CmdMatches {
	iter := NewArgIterator(args)
	return c.matches(iter)
}

func (c *Command) matches(iter *ArgIterator) *CmdMatches {
	matches := NewCmdMatches(iter.Next())

	if iter.NextIsFlag() {
		for _, flag := range c.Flags {
			log.Printf("checking %+v\n", *flag)
			if match := flag.Matches(iter); match != nil {
				matches.Flags[flag.Name] = *match
				break
			}
		}
	} else {
		// next is subcommand
		for _, sub := range c.Subcommands {
			log.Printf("checking subcommand %s for %s\n", sub.Name, iter.Peek())
			if sub.Name == iter.Peek() {
				matches.Subcommand = sub.matches(iter)
				break
			}
		}
	}

	return matches
}

type CmdMatches struct {
	Name       string
	Flags      map[string]FlagMatch
	Subcommand *CmdMatches
}

func NewCmdMatches(name string) *CmdMatches {
	return &CmdMatches{
		Name:  name,
		Flags: make(map[string]FlagMatch),
	}
}

func (c *CmdMatches) HasFlag(name string) bool {
	_, ok := c.Flags[name]
	return ok
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
