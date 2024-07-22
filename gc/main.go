package main

import (
	"fmt"
	"net/http"
)

const SIZE = 1024 * 1024

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// var byteArray = make([]byte, SIZE)
	// fmt.Fprintf(w, "hello world %d", len(byteArray))

	fmt.Fprintf(w, "hello world")
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}

// go test -bench=. -benchmem -benchtime=10s -parallel=4 -v
