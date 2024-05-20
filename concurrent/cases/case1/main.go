package main

import "time"

type Token struct{}

func newWorker(id int, ch chan Token, nextCh chan Token) {
	for {
		token := <-ch
		println(id)
		time.Sleep(time.Second)
		nextCh <- token
	}
}

func main() {
	chs := []chan Token{make(chan Token), make(chan Token), make(chan Token), make(chan Token)}
	for i := range 4 {
		go newWorker(i, chs[i], chs[(i+1)%4])
	}

	chs[0] <- Token{}
	select {}
}
