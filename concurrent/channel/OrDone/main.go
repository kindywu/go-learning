package main

import (
	"fmt"
	"reflect"
	"time"
)

func main() {
	start := time.Now()
	ch := or(sig(10*time.Second),
		sig(5*time.Second),
		sig(20*time.Second),
		sig(3*time.Second),
		sig(50*time.Second),
		sig(7*time.Second),
		sig(4*time.Second),
		sig(6*time.Second))
	if ch != nil {
		<-ch
	}
	fmt.Printf("done after %v\n", time.Since(start))
	start = time.Now()
	ch = or()
	if ch != nil {
		<-ch
	}
	fmt.Printf("done after %v\n", time.Since(start))
	start = time.Now()
	ch2 := or2(sig(10*time.Second),
		sig(5*time.Second),
		sig(20*time.Second),
		sig(4*time.Second),
		sig(6*time.Second))
	if ch2 != nil {
		<-ch2
	}
	fmt.Printf("done after %v\n", time.Since(start))
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func or(chs ...<-chan interface{}) <-chan interface{} {
	switch len(chs) {
	case 0:
		return nil
	case 1:
		return chs[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(chs) {
		case 2:
			select {
			case <-chs[0]:
			case <-chs[1]:
			}
		default:
			{
				m := len(chs) / 2
				select {
				case <-or(chs[:m]...):
				case <-or(chs[m:]...):
				}
			}
		}
	}()
	return orDone
}

func or2(chs ...<-chan interface{}) <-chan interface{} {
	switch len(chs) {
	case 0:
		return nil
	case 1:
		return chs[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		var cases []reflect.SelectCase
		for _, c := range chs {
			cases = append(cases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(c),
			})
		}
		reflect.Select(cases)
	}()
	return orDone
}
