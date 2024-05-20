package queue

import (
	"log/slog"
	"sync"
	"sync/atomic"
	"testing"
)

var logger = slog.With("queue", "1.0")

func TestQueue(t *testing.T) {
	// slog.SetLogLoggerLevel(slog.LevelInfo)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	turn := 10
	round := 50
	q := NewLKQueue()
	wg := sync.WaitGroup{}
	wgCount := turn * round
	logger.Debug("WaitGroup", "count", wgCount)
	wg.Add(wgCount)
	var enqueue int32 = 0
	var dequeue int32 = 0
	for range turn {
		go func() {
			for i := range round {
				logger.Debug("enqueue-test", "count", atomic.AddInt32(&enqueue, 1), "value", i)
				q.Enqueue(i)
			}
		}()
		go func() {
			for {
				v := q.Dequeue()
				if v != nil {
					logger.Debug("dequeue-test", "count", atomic.AddInt32(&dequeue, 1), "value", v)
					wg.Done()
				}
			}
		}()
	}
	wg.Wait()
}
