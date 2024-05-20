package main

import (
	"sync"
	"time"
)

// 使用chan实现互斥锁
type Mutex struct {
	ch chan bool
}

// 使用锁需要初始化
func NewMutex() *Mutex {
	mu := &Mutex{make(chan bool, 1)}
	mu.ch <- true
	return mu
}

// 请求锁，直到获取到
func (m *Mutex) Lock() {
	<-m.ch
}

// 获取锁，超时放弃
func (m *Mutex) LockTimeout(timeout time.Duration) bool {
	timer := time.NewTimer(timeout)
	select {
	case <-m.ch:
		timer.Stop()
		return true
	case <-timer.C:
	}
	return false
}

// 解锁
func (m *Mutex) Unlock() {
	select {
	case m.ch <- true:
	default:
		panic("unlock of unlocked mutex")
	}
}

// 尝试获取锁
func (m *Mutex) TryLock() bool {
	select {
	case <-m.ch:
		return true
	default:
		return false
	}
}

// 锁是否已被持有
func (m *Mutex) IsLocked() bool {
	return len(m.ch) == 0
}

func main() {
	count := 10000
	wg := sync.WaitGroup{}
	wg.Add(count)
	m := NewMutex()
	num := 0
	for range count {
		go func() {
			defer wg.Done()
			m.Lock()
			defer m.Unlock()
			num++
		}()

		go func() {
			if got := m.TryLock(); got {
				println("I got the lock", m.IsLocked())
				defer m.Unlock()
			}
		}()
	}
	wg.Wait()
	println(num)
}
