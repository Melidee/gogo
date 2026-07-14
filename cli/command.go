package cli

import (
	"fmt"
	"os"
	"strings"
)

type Cmd interface {
	GetName() string
	GetAuthor() string
	GetVersion() string
	GetAbout() string
	Apply(args []string) error
	apply(iter *argIter) error
	applyFlag(iter *argIter) error
	nextSubcmd(iter *argIter) Cmd
}

type Empty struct{}

type Command[T any] struct {
	name        string
	author      string
	version     string
	about       string
	usage       string
	customHelp  func(Command[T]) string
	flags       []*Flag[T]
	subcommands []Cmd
	state       T
	action      func(ctx Context[T], value string) error
}

func NewCommand[T any](name string, init T) *Command[T] {
	return &Command[T]{
		name:   name,
		state:  init,
		action: func(ctx Context[T], value string) error { return nil },
	}
}

func SimpleCommand(name string) *Command[Empty] {
	return NewCommand(name, Empty{})
}

func (c *Command[T]) init() {
	if c.customHelp == nil {
		c.addHelpCommands()
	}
}

func (c *Command[T]) addHelpCommands() {
	helpCmd := NewCommand("help", Empty{}).
		About("Print help message and exit").
		Action(func(ctx Context[Empty], value string) error {
			c.PrintHelp()
			os.Exit(0)
			return nil
		})
	c.subcommands = append(c.subcommands, helpCmd)
	flag := NewFlag[T]("help").
		Short('h').
		Long("help").
		About("Print help message and exit").
		Action(func(ctx Context[T], value string) error {
			c.PrintHelp()
			os.Exit(0)
			return nil
		})
	c.flags = append(c.flags, flag)
}

/* --- Getters --- */

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

/* Help formatting */

func (c *Command[T]) GetHelpString() string {
	help := Style(c.name, Blue, Bold)
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
		help += "\n" + Style("Usage", Green) + ":\n"
		help += "    " + c.GetUsage()
	}
	help += c.formatCommandsHelp()
	help += c.formatOptionsHelp()
	return help
}

func (c *Command[T]) Help() *Command[T] {

	return c
}

func (c *Command[T]) formatCommandsHelp() string {
	if len(c.subcommands) == 0 {
		return ""
	}

	longPadSize := 0
	for _, cmd := range c.subcommands {
		if len(cmd.GetName()) > longPadSize {
			longPadSize = len(cmd.GetName())
		}
	}
	longPadSize += 2

	help := "\n\n" + Style("Commands", Green) + ":"
	for _, cmd := range c.subcommands {
		help += fmt.Sprintf("\n    %s", Style(cmd.GetName(), Blue))
		pad := strings.Repeat(" ", longPadSize-len(cmd.GetName()))
		help += pad + cmd.GetAbout()
	}
	return help
}

func (c *Command[T]) formatOptionsHelp() string {
	// format short and long flags
	var flagStrings []string 
	for _, flag := range c.flags {
		s := "    "
		if flag.short != 0 {
			short := fmt.Sprintf("-%c", flag.short)
			s = Style(short, Blue) + ", "
		}
		s += Style("--" + flag.long, Blue)
		if flag.argName != "" {
			s += " <" + flag.argName + ">"
		}
		s = "    " + s + "  " // add padding
		flagStrings = append(flagStrings, s)
	}

	// find the longest flag string to calculate how much padding to add
	longestFlag := 0
	for _, flagStr := range flagStrings {
		if longestFlag < len(flagStr) {
			longestFlag = len(flagStr)
		}
	}
	
	// add padding and about text to flags
	for i, flag := range c.flags {
		flagStr := flagStrings[i]
		padSize := longestFlag - len(flagStr)
		pad := strings.Repeat(" ", padSize)
		flagStrings[i] = flagStr + pad + flag.about
	}

	header := "\n\n" + Style("Options", Green) + ":\n"
	return header + strings.Join(flagStrings, "\n")
}

func (c *Command[T]) PrintHelp() {
	fmt.Println(c.GetHelpString())
}

/* Builder methods */

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
		About("Print version info and exit").Action(func(ctx Context[Empty], value string) error {
		fmt.Printf("%s version v%s\n", c.name, c.version)
		os.Exit(0)
		return nil
	})
	c.subcommands = append(c.subcommands, subCmd)

	flag := NewFlag[T]("version").
		Short('V').
		Long("version").
		About("Print version info and exit").
		Action(func(ctx Context[T], value string) error {
			fmt.Printf("%s version v%s\n", ctx.cmd.name, ctx.cmd.version)
			os.Exit(0)
			return nil
		})
	c.flags = append(c.flags, flag)
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

func (c *Command[T]) Action(action func(ctx Context[T], value string) error) *Command[T] {
	c.action = action
	return c
}

/* Action running */

func (c *Command[T]) Apply(args []string) error {
	c.init()
	return c.apply(newArgsIter(args))
}

func (c *Command[T]) apply(iter *argIter) error {
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
	return c.action(NewContext(c), cmdArg)
}

func (c *Command[T]) applyFlag(iter *argIter) error {
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

	return flag.action(NewContext(c), value)
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
