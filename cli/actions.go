package cli

import "strconv"

func ActionSet[T any](dest func(state T) *string) func(ctx Context[T], value string) error {
	return func(ctx Context[T], value string) error {
		*dest(*ctx.State()) = value
		return nil
	}
}

func ActionSetInt[T any](dest func(state T) *int) func(ctx Context[T], value string) error {
	return func(ctx Context[T], value string) error {
		num, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		*dest(*ctx.State()) = num
		return nil
	}
}

func ActionSetFloat32[T any](dest func(state T) *float32) func(ctx Context[T], value string) error {
	return func(ctx Context[T], value string) error {
		num, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		*dest(*ctx.State()) = float32(num)
		return nil
	}
}

func ActionSetFloat64[T any](dest func(state T) *float64) func(ctx Context[T], value string) error {
	return func(ctx Context[T], value string) error {
		num, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		*dest(*ctx.State()) = num
		return nil
	}
}

func ActionSetTrue[T any](dest func(state T) *bool) func(ctx Context[T], value string) error {
	return func(ctx Context[T], value string) error {
		*dest(*ctx.State()) = true
		return nil
	}
}

func ActionSetFalse[T any](dest func(state T) *bool) func(ctx Context[T], value string) error {
	return func(ctx Context[T], value string) error {
		*dest(*ctx.State()) = false
		return nil
	}
}

func ActionToggle[T any](dest func(state T) *bool) func(ctx Context[T], value string) error {
	return func(ctx Context[T], value string) error {
		*dest(*ctx.State()) = !*dest(*ctx.State())
		return nil
	}
}