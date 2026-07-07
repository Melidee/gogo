package hopeful

import (
	"fmt"
)

func test() {

}

type Cmd interface {
	Apply(args []string)
}

func NewCmd(name string) *Command[struct{}] {
	return &Command[struct{}]{
		Name: name,
	}
}

type Command[T any] struct {
	Name        string
	Author      string
	Version     string
	About       string
	Help        string
	Flags       []*Flag[T]
	Subcommands []Cmd
	State       T
	Action      func(ctx Context[T], value string)
}

func NewCommand[T any](name string, init T) *Command[T] {
	return &Command[T]{
		Name:  name,
		State: init,
	}
}

func NewCommand_(name string) *Command[struct{}] {
	return &Command[struct{}]{
		Name: name,
	}
}

func (c *Command[T]) SetName(name string) *Command[T] {
	c.Name = name
	return c
}

func (c *Command[T]) SetAuthor(author string) *Command[T] {
	c.Author = author
	return c
}

func (c *Command[T]) versionCommands() (*Command[struct{}], *Flag[T]) {
	subCmd := NewCommand_("version").SetHelp("Print the version of the program.")
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

func (c *Command[T]) CallAction(value string) *Command[T] {
	if c.Action != nil {
		c.Action(NewContext(c), value)
	}
	return c
}

func (c *Command[T]) Apply(args []string) {
	iter := NewArgsIter(args)
	iter.Next()
	c.apply(iter)
}

func (c *Command[T]) apply(iter *ArgIter) {
	cmdArg := ""
	for iter.HasNext() {
		if iter.NextIsFlag() {
			c.applyFlag(iter)
		} else if c.NextIsSubcmd(iter) {
			
		} else {
			cmdArg = iter.Next()
			break
		}
	}
	c.CallAction(cmdArg)
}

func (c *Command[T]) applyFlag(iter *ArgIter) {
	// get the token for the flag we're looking at
	flagToken := iter.Next()
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

	value := ""
	if flag.TakesValue && iter.NextIsValue() {
		value = iter.Next()
	}

	flag.action(NewContext(c), value)
}

func (c *Command[T]) NextIsSubcmd(iter *ArgIter) bool {
	if !iter.HasNext() {
		return false
	}
	subcmdToken := iter.Peek()
	for _, sub := range c.Subcommands {
		if sub.(*Command[any]).Name == subcmdToken {
			return true
		}
	}
	return false
}

func (c *Command[T]) ToCmd() Cmd {
	return c
}
