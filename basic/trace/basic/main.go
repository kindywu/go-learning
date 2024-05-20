package main

func Trace(name string) func() {
	println("enter", name)
	return func() {
		println("exit", name)
	}
}

func foo() {
	defer Trace("foo")()
	bar()
}

func bar() {
	defer Trace("bar")()
}

func main() {
	defer Trace("main")()
	foo()
}
