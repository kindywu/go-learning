package main

import "sync"

type Queue[T any] struct {
	data []T
	mu   sync.Mutex
}

func NewQueue[T any](n int) *Queue[T] {
	return &Queue[T]{data: make([]T, 0, n)}
}

func (q *Queue[T]) Enqueue(v T) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, v)
}

func (q *Queue[T]) Dequeue() *T {
	q.mu.Lock()
	if len(q.data) == 0 {
		q.mu.Unlock()
		return nil
	}
	defer q.mu.Unlock()
	v := q.data[0]
	q.data = q.data[1:]
	return &v
}

func main() {
	queue := NewQueue[int](10)
	queue.Enqueue(10)
	v := queue.Dequeue()
	println(*v)
}
