package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -lhello
#include "hello.h"
*/
import "C"

import (
	"fmt"
	"os"
)

func main() {
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("无法获取当前目录:", err)
	} else {
		// 打印当前工作目录
		fmt.Println("当前执行目录:", dir)
	}
	C.SayHello()
	fmt.Println("Success!")
}

//gcc -fPIC -c hello.c
//gcc -shared -o libhello.so hello.o
//export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/workspaces/Demo/basic/c/gcc
