package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

var goroutineSpace = []byte("goroutine ")

func curGoroutineID() uint64 {
	b := make([]byte, 64)
	len := runtime.Stack(b, false)
	b = b[:len]
	// fmt.Println("b", string(b))
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		panic(fmt.Sprintf("No space found in %q", b))
	}
	b = b[:i]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("Fail to parse goroutine ID out of %q:%v", b, err))
	}
	return n
}

func printTrace(id uint64, name, arrow string, indent int) {
	indents := ""
	for range indent {
		indents += "  "
	}
	fmt.Printf("g[%05d]:%s%s%s\n", id, indents, arrow, name)
}

var mu sync.Mutex
var m = make(map[uint64]int)

func Trace() func() {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("Can not found caller")
	}
	fn := runtime.FuncForPC(pc)
	name := fn.Name()
	gid := curGoroutineID()

	mu.Lock()
	defer mu.Unlock()

	indent := m[gid]
	m[gid] = indent + 1

	printTrace(gid, name, "->", indent+1)
	return func() {
		indent := m[gid]
		m[gid] = indent - 1
		printTrace(gid, name, "->", indent)
	}
}

func A1() {
	defer Trace()()
	B1()
}

func B1() {
	defer Trace()()
	C1()
}

func C1() {
	defer Trace()()
	D()
}

func D() {
	defer Trace()()
}

func A2() {
	defer Trace()()
	B2()
}

func B2() {
	defer Trace()()
	C1()
}

func C2() {
	defer Trace()()
	D()
}

func main() {
	defer Trace()()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		A2()
	}()

	A1()
	wg.Wait()
}
