package main

import "sync"

type Pool struct {
	ch chan []byte
}

func new(size int, chanLen int) *Pool {
	ch := make(chan []byte, chanLen)
	for n := 0; n < chanLen; n++ {
		ch <- make([]byte, size)
	}
	return &Pool{
		ch: ch,
	}
}

func (p *Pool) Get() *[]byte {
	buf := <-p.ch
	return &buf
}

func (p *Pool) Put(buf *[]byte) {
	p.ch <- *buf
}

var once sync.Once
var myPool *Pool

func NewPool(size int, chanLen int) *Pool {
	once.Do(func() {
		myPool = new(size, chanLen)
		// 初始化singleton的逻辑
	})
	return myPool
}
