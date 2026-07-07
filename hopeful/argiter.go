package hopeful

type ArgIter struct {
	args []string
	cursor int
}

func NewArgsIter(args []string) *ArgIter {
	return &ArgIter{
		args: args,
		cursor: 0,
	}
}

func (a *ArgIter) Peek() string {
	if a.cursor >= len(a.args) {
		return ""
	}
	return  a.args[a.cursor]
}

func (a *ArgIter) Next() string {
	next := a.Peek()
	if next != "" {
		a.cursor++
	}
	return next
}

func (a *ArgIter) HasNext() bool {
	return a.cursor < len(a.args)
}

func (a *ArgIter) NextIsFlag() bool {
	return a.Peek() != "" && a.Peek()[0] == '-'
}

func (a *ArgIter) NextIsValue() bool {
	return a.Peek() != "" && !a.NextIsFlag()
}