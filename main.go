package main

import (
	"runtime"

	_ "go.uber.org/automaxprocs"
)

var ()

func main() {
	println(runtime.NumCPU())
	println(runtime.GOMAXPROCS(0))
}
