package main

/*
#include "hello.c"
*/
import "C"
import (
	"fmt"
)

func main() {
	ret := C.SayHello()
	fmt.Println(ret)
	sum := C.Add(C.int(3), C.int(4))
	fmt.Println("The sum is", sum)
}
