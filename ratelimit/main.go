package main

import (
	"log"

	"go.uber.org/ratelimit"
)

func main() {
	r1 := ratelimit.New(1, ratelimit.WithSlack(3))

	for i := 0; i < 10; i++ {
		log.Println("got #", i, r1.Take())
	}
}
