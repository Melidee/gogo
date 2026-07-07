package hopeful

import "strconv"

type Flag[T any] struct {
	Name       string
	Short      rune
	Long       string
	About      string
	TakesValue bool
	HasDefault bool
	Default    string
	Required   bool
	action     func(ctx Context[T], value string)
}

func NewFlag[T any](name string) *Flag[T] {
	return &Flag[T]{
		Name: name,
	}
}

func NewFlag_(name string) *Flag[struct{}] {
	return &Flag[struct{}]{
		Name: name,
	}
}

func (f *Flag[T]) SetShort(short rune) *Flag[T] {
	f.Short = short
	return f
}

func (f *Flag[T]) SetLong(long string) *Flag[T] {
	f.Long = long
	return f
}

func (f *Flag[T]) SetAbout(about string) *Flag[T] {
	f.About = about
	return f
}

func (f *Flag[T]) SetDefault(defaultValue string) *Flag[T] {
	f.HasDefault = true
	f.Default = defaultValue
	return f
}

func (f *Flag[T]) SetRequired(yes bool) *Flag[T] {
	f.Required = yes
	return f
}

func (f *Flag[T]) Action(action func(ctx Context[T], value string)) *Flag[T] {
	f.action = action
	return f
}

func (f *Flag[T]) CallAction(ctx Context[T], value string)
