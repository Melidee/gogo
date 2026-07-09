package cli

import (
	"context"
	"time"
)

type Context[T any] struct {
	cmd *Command[T]
}

func NewContext[T any](cmd *Command[T]) Context[T] {
	return Context[T]{
		cmd: cmd,
	}
}

func (c Context[T]) State() *T {
	return &c.cmd.state
}

func (c Context[T]) Flags() []*Flag[T] {
	return c.cmd.flags
}

func (c Context[T]) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (c Context[T]) Done() <-chan struct{} {
	return nil
}

func (c Context[T]) Err() error {
	return nil
}

func (c Context[T]) Value(key any) any {
	return nil
}

type a context.Context