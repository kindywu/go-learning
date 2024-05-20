package common

import "context"

type Runtime struct {
	Add func(ctx context.Context, x int, y int) (int, error)
}

func (ri *Runtime) RegisterAdd(fn func(ctx context.Context, x int, y int) (int, error)) error {
	ri.Add = fn
	return nil
}
