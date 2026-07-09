package cli

type argIter struct {
	args []string
	cursor int
}

func newArgsIter(args []string) *argIter {
	return &argIter{
		args: args,
		cursor: 0,
	}
}

func (a *argIter) peek() string {
	if !a.hasNext() {
		return ""
	}
	return  a.args[a.cursor]
}

func (a *argIter) next() string {
	next := a.peek()
	if next != "" {
		a.cursor++
	}
	return next
}

func (a *argIter) hasNext() bool {
	return a.cursor <= len(a.args)
}

func (a *argIter) nextIsFlag() bool {
	return a.peek() != "" && a.peek()[0] == '-'
}

func (a *argIter) nextIsValue() bool {
	return a.peek() != "" && !a.nextIsFlag()
}