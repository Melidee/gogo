package cli

import (
	"fmt"
)

type Cmd interface {
	GetName() string
	GetAuthor() string
	GetVersion() string
	GetAbout() string
	GetHelp() string
	Apply(args []string)
	apply(iter *argIter)
	applyFlag(iter *argIter)
	nextSubcmd(iter *argIter) Cmd
}

type Empty struct{}

type Command[T any] struct {
	Name        string
	Author      string
	Version     string
	About       string
	Help        string
	Flags       []*Flag[T]
	Subcommands []Cmd
	State       T
	action      func(ctx Context[T], value string)
}

func NewCommand[T any](name string, init T) *Command[T] {
	return &Command[T]{
		Name:  name,
		State: init,
	}
}

func (c *Command[T]) GetName() string {
	return c.Name
}

func (c *Command[T]) GetAuthor() string {
	return c.Author
}

func (c *Command[T]) GetVersion() string {
	return c.Version
}

func (c *Command[T]) GetAbout() string {
	return c.About
}

func (c *Command[T]) GetHelp() string {
	return c.Help
}

func (c *Command[T]) SetName(name string) *Command[T] {
	c.Name = name
	return c
}

func (c *Command[T]) SetAuthor(author string) *Command[T] {
	c.Author = author
	return c
}

func (c *Command[T]) versionCommands() (*Command[Empty], *Flag[T]) {
	subCmd := NewCommand[Empty]("version", Empty{}).SetHelp("Print the version of the program.")
	flag := NewFlag[T]("version").SetShort('V').SetLong("version").Action(func(ctx Context[T], value string) {
		fmt.Printf("%s version v%s", ctx.cmd.Name, ctx.cmd.Version)
	})
	return subCmd, flag
}

func (c *Command[T]) SetVersion(version string) *Command[T] {
	c.Version = version
	subCmd, flag := c.versionCommands()
	c.Subcommands = append(c.Subcommands, subCmd)
	c.Flags = append(c.Flags, flag)
	return c
}

func (c *Command[T]) SetHelp(help string) *Command[T] {
	c.Help = help
	return c
}

func (c *Command[T]) AddFlag(flag *Flag[T]) *Command[T] {
	c.Flags = append(c.Flags, flag)
	return c
}

func (c *Command[T]) AddSubcommand(cmd Cmd) *Command[T] {
	c.Subcommands = append(c.Subcommands, cmd)
	return c
}

func (c *Command[T]) Action(action func(ctx Context[T], value string)) *Command[T] {
	c.action = action
	return c
}

func (c *Command[T]) Apply(args []string) {
	c.apply(newArgsIter(args))
}

func (c *Command[T]) apply(iter *argIter) {
	iter.next() // consume the command name
	cmdArg := ""
	for iter.hasNext() {
		if iter.nextIsFlag() {
			c.applyFlag(iter)
		} else if subcmd := c.nextSubcmd(iter); subcmd != nil {
			subcmd.apply(iter)
		} else {
			cmdArg = iter.next()
			break
		}
	}
	c.action(NewContext(c), cmdArg)
}

func (c *Command[T]) applyFlag(iter *argIter) {
	// get the token for the flag we're looking at
	flagToken := iter.next()
	if len(flagToken) < 2 {
		panic("single dash")
	}

	// find the matching flag
	var flag *Flag[T]
	for _, f := range c.Flags {
		matchesShort := flagToken[0] == '-' && flagToken[1] == byte(f.Short)
		matchesLong := flagToken[0:2] == "--" && flagToken[2:] == f.Long
		if matchesShort || matchesLong {
			flag = f
			break
		}
	}
	if flag == nil {
		panic(fmt.Sprintf("unknown flag: %s", flagToken))
	}

	// capture value if the flag takes one
	value := ""
	if flag.TakesValue && iter.nextIsValue() {
		value = iter.next()
	}

	flag.action(NewContext(c), value)
}

func (c *Command[T]) nextSubcmd(iter *argIter) Cmd {
	if !iter.hasNext() {
		return nil
	}
	subcmdToken := iter.peek()
	for _, sub := range c.Subcommands {
		if sub.GetName() == subcmdToken {
			return sub
		}
	}
	return nil
}

func (c *Command[T]) ToCmd() Cmd {
	return c
}
