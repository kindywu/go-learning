package main

import (
	"context"
	"fmt"
	"log"
	"main/plugin/common"
	"os"
	"plugin"
)

func main() {
	path := "./plugin/sample.so"
	info, err := os.Stat(path)
	if err != nil {
		log.Fatalln("Go module is not exist", path, err)
	}

	name := info.Name()
	fmt.Println("Go module name is", name)

	p, err := plugin.Open(path)
	if err != nil {
		log.Fatalln("Could not open Go module", path, err)
	}

	// Look up the required initialisation function.
	f, err := p.Lookup("InitModule")
	if err != nil {
		log.Fatalln("Error looking up InitModule function in Go module", name)
	}

	ri := common.Runtime{}
	fn, ok := f.(func(string, *common.Runtime) error)
	if !ok {
		log.Fatal("Error reading InitModule function in Go module", name)
	}
	fn("kindy", &ri)

	if ri.Add != nil {
		println(ri.Add(context.Background(), 1, 100))
	}
}
