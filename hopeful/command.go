package hopeful

import "fmt"

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
	Action      func(cmd Command[T], flags []*Flag[T], state *T)
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

func (c *Command[T]) CallAction() *Command[T] {
	if c.Action != nil {
		c.Action(*c, c.Flags, &c.State)
	}
	return c
}

func (c *Command[T]) Apply([]string) {

}

func (c *Command[T]) ToCmd() Cmd {
	return c
}
