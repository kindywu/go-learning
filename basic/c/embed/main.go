package main

/*
#include <stdio.h>
int SayHello() {
 puts("Hello World");
    return 0;
}
int Add(int a, int b) {
    return a + b;
}
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
