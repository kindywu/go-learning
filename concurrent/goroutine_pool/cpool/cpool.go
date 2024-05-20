package cpool

import (
	"fmt"
	"log/slog"
	"sync"
)

var logger = slog.With("cpool", "1.0")

type TaskResult struct {
	Err    error
	Result interface{}
}

type Task struct {
	Name string
	Do   func() (interface{}, error)
	Done chan TaskResult
}

func NewTask(name string, do func() (interface{}, error)) *Task {
	return &Task{
		Name: name,
		Do:   do,
		Done: make(chan TaskResult, 1),
	}
}

type Pool struct {
	wg   *sync.WaitGroup
	task chan *Task
}

func NewPool(size int) *Pool {
	logger.Debug("new pool", "size", size)
	p := &Pool{
		wg:   new(sync.WaitGroup),
		task: make(chan *Task),
	}
	p.wg.Add(size)
	for i := 0; i < size; i++ {
		go p.worker(i)
	}
	return p
}

func (p *Pool) Submit(task *Task) {
	if task != nil {
		logger.Debug("submit", "task", task.Name)
		p.task <- task
	}
}

func (p *Pool) Stop() {
	logger.Debug("stop")
	close(p.task)
	p.wg.Wait()
}

func (p *Pool) worker(num int) {
	defer p.wg.Done()

	logger.Debug("enter", "worker", num)

	for t := range p.task {
		p.exectue(num, t)
	}

	logger.Debug("exit", "worker", num)
}

func (p *Pool) exectue(num int, t *Task) {
	defer func() {
		if r := recover(); r != nil {
			logger.Warn("panic and recover", "worker", num, "panic", r)
			t.Done <- TaskResult{Result: nil, Err: fmt.Errorf("%+v", r)}
			close(t.Done)
		}
	}()

	logger.Debug("task executing", "worker", num, "task", t.Name)
	result, err := t.Do()
	logger.Debug("task done", "worker", num, "task", t.Name)
	t.Done <- TaskResult{Result: result, Err: err}
	close(t.Done)
}
