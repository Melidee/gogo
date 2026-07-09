package cli

import "strconv"

type Flag[T any] struct {
	name       string
	short      rune
	long       string
	about      string
	takesValue bool
	hasDefault bool
	default_   string
	required   bool
	argName    string
	action     func(ctx Context[T], value string) error
}

func NewFlag[T any](name string) *Flag[T] {
	return &Flag[T]{
		name: name,
	}
}

func (f *Flag[T]) Short(short rune) *Flag[T] {
	f.short = short
	return f
}

func (f *Flag[T]) Long(long string) *Flag[T] {
	f.long = long
	return f
}

func (f *Flag[T]) About(about string) *Flag[T] {
	f.about = about
	return f
}

func (f *Flag[T]) Default(defaultValue string) *Flag[T] {
	f.hasDefault = true
	f.default_ = defaultValue
	return f
}

func (f *Flag[T]) Required(yes bool) *Flag[T] {
	f.required = yes
	return f
}

func (f *Flag[T]) ArgName(argName string) *Flag[T] {
	f.argName = argName
	return f
}

func (f *Flag[T]) Action(action func(ctx Context[T], value string) error) *Flag[T] {
	f.action = action
	return f
}

func (f *Flag[T]) ActionSet(dest func(state T) *string) *Flag[T] {
	f.action = func(ctx Context[T], value string) error {
		*dest(*ctx.State()) = value
		return nil
	}
	return f
}

func (f *Flag[T]) ActionSetInt(dest func(state T) *int) *Flag[T] {
	f.action = func(ctx Context[T], value string) error {
		num, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*dest(*ctx.State()) = num
		return nil
	}
	return f
}

func (f *Flag[T]) ActionSetTrue(dest func(state T) *bool) *Flag[T] {
	f.action = func(ctx Context[T], value string) error {
		*dest(*ctx.State()) = true
		return nil
	}
	return f
}

func (f *Flag[T]) ActionSetFalse(dest func(state T) *bool) *Flag[T] {
	f.action = func(ctx Context[T], value string) error {
		*dest(*ctx.State()) = false
		return nil
	}
	return f
}
