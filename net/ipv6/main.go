package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})
	// err := http.ListenAndServe("[::]:8083", nil)
	err := http.ListenAndServe("[0:0:0:0:0:0:0:1]:8083", nil)
	if err != nil {
		println(err.Error())
	}
}
