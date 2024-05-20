package main

import (
	"context"
	"fmt"
	"main/plugin/common"
)

func InitModule(name string, ri *common.Runtime) error {
	fmt.Println("hello", name)
	ri.RegisterAdd(Add)
	return nil
}

func Add(ctx context.Context, x int, y int) (int, error) {
	return x + y, nil
}
