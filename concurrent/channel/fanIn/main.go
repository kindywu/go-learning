package main

import (
	"bufio"
	"os"
	"reflect"
	"time"
)

func fanIn(chans ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		var cases []reflect.SelectCase
		for _, c := range chans {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}

		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok {
				cases = append(cases[:i], cases[i+1:]...)
				continue
			}
			out <- v.Interface()
		}
	}()
	return out
}

func main() {
	ch1 := make(chan interface{})
	ch2 := make(chan interface{})

	go func() {
		for {
			ch1 <- 1
			ch2 <- 2
			time.Sleep(time.Second)
		}
	}()

	go func() {
		out := fanIn(ch1, ch2)
		for v := range out {
			println(v.(int))
		}
	}()

	bufio.NewReader(os.Stdin).ReadString('\n')
}
