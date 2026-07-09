package cli

import (
	"fmt"
	"os"
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
	name    string
	author  string
	version string
	about   string
	usage   string
	help    string

	flags       []*Flag[T]
	subcommands []Cmd
	state       T
	action      func(ctx Context[T], value string)
}

func NewCommand[T any](name string, init T) *Command[T] {
	return &Command[T]{
		name:   name,
		state:  init,
		action: func(ctx Context[T], value string) {},
	}
}

func (c *Command[T]) GetName() string {
	return c.name
}

func (c *Command[T]) GetAuthor() string {
	return c.author
}

func (c *Command[T]) GetVersion() string {
	return c.version
}

func (c *Command[T]) GetAbout() string {
	return c.about
}

func (c *Command[T]) GetUsage() string {
	return c.usage
}

func (c *Command[T]) GetHelp() string {
	return c.help
}

func (c *Command[T]) GetHelpString() string {
	help := c.name
	if c.GetVersion() != "" {
		help += fmt.Sprintf(" v%s", c.GetVersion())
	}
	help += "\n"
	if c.GetAuthor() != "" {
		help += fmt.Sprintf("%s\n", c.GetAuthor())
	}
	if c.GetAbout() != "" {
		help += fmt.Sprintf("%s\n", c.GetAbout())
	}
	if c.GetUsage() != "" {
		help += fmt.Sprintf("\nUsage:\n    %s\n", c.GetUsage())
	}
	help += c.formatCommandsHelp()
	help += c.formatOptionsHelp()
	return help
}

func (c *Command[T]) formatCommandsHelp() string {
	if len(c.subcommands) == 0 {
		return ""
	}
	help := "\nCommands:\n"
	for _, cmd := range c.subcommands {
		help += fmt.Sprintf("    %s\n", cmd.GetName())
	}
	return help
}

func (c *Command[T]) formatOptionsHelp() string {
	if len(c.flags) == 0 {
		return ""
	}
	help := "\nOptions:"
	for _, flag := range c.flags {
		short := "    "
		if flag.short != 0 {
			short = fmt.Sprintf("-%c, ", flag.short)
		}
		help += fmt.Sprintf("\n    %s--%s\t%s", short, flag.long, flag.about)
	}
	return help
}

func (c *Command[T]) PrintHelp() {
	fmt.Println(c.GetHelpString())
}

func (c *Command[T]) Name(name string) *Command[T] {
	c.name = name
	return c
}

func (c *Command[T]) Author(author string) *Command[T] {
	c.author = author
	return c
}

func (c *Command[T]) Version(version string) *Command[T] {
	c.version = version

	subCmd := NewCommand("version", Empty{}).
		About("Print version info and exit").Action(func(ctx Context[Empty], value string) {
		fmt.Printf("%s version v%s\n", c.name, c.version)
		os.Exit(0)
	})
	c.subcommands = append(c.subcommands, subCmd)

	flag := NewFlag[T]("version").
		Short('V').
		Long("version").
		About("Print version info and exit").
		Action(func(ctx Context[T], value string) {
			fmt.Printf("%s version v%s\n", ctx.cmd.name, ctx.cmd.version)
			os.Exit(0)
		})
	c.flags = append(c.flags, flag)
	return c
}

func (c *Command[T]) Help(help string) *Command[T] {
	helpCmd := NewCommand("help", Empty{}).
		About("Print help message and exit").
		Action(func(ctx Context[Empty], value string) {
			c.PrintHelp()
			os.Exit(0)
		})
	c.subcommands = append(c.subcommands, helpCmd)
	flag := NewFlag[T]("help").
		Short('h').
		Long("help").
		About("Print help message and exit").
		Action(func(ctx Context[T], value string) {
			c.PrintHelp()
			os.Exit(0)
		})
	c.flags = append(c.flags, flag)
	c.help = help
	return c
}

func (c *Command[T]) About(about string) *Command[T] {
	c.about = about
	return c
}

func (c *Command[T]) Usage(usage string) *Command[T] {
	c.usage = usage
	return c
}

func (c *Command[T]) Flag(flag *Flag[T]) *Command[T] {
	c.flags = append(c.flags, flag)
	return c
}

func (c *Command[T]) Subcommand(cmd Cmd) *Command[T] {
	c.subcommands = append(c.subcommands, cmd)
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
	for _, f := range c.flags {
		matchesShort := flagToken[0] == '-' && flagToken[1] == byte(f.short)
		matchesLong := flagToken[0:2] == "--" && flagToken[2:] == f.long
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
	if flag.takesValue && iter.nextIsValue() {
		value = iter.next()
	}

	flag.action(NewContext(c), value)
}

func (c *Command[T]) nextSubcmd(iter *argIter) Cmd {
	if !iter.hasNext() {
		return nil
	}
	subcmdToken := iter.peek()
	for _, sub := range c.subcommands {
		if sub.GetName() == subcmdToken {
			return sub
		}
	}
	return nil
}
