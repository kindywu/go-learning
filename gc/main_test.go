package main

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -bench=BenchmarkIndexHandler1 -benchmem -benchtime=10s -v
// go test -bench=BenchmarkIndexHandler2 -benchmem -benchtime=10s -v
// go test -bench=BenchmarkIndexHandler3 -benchmem -benchtime=10s -v

// go test -bench . -benchmem -benchtime=10s -parallel=4 -v

func TestPool(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				var byteArray = pool.Get().(*[]byte)
				size := len(*byteArray)
				pool.Put(byteArray)
				assert.Equal(t, SIZE, size)
			}
		}()
	}
	wg.Wait()
}

func BenchmarkIndexHandler1(b *testing.B) {
	// 创建请求
	req, _ := http.NewRequest("GET", "/", nil)

	b.RunParallel(func(p *testing.PB) {
		w := httptest.NewRecorder()
		for p.Next() {
			w.Body.Reset()
			indexHandler1(w, req)
		}
	})

	printMemStats(b)
}

func BenchmarkIndexHandler2(b *testing.B) {
	// 创建请求
	req, _ := http.NewRequest("GET", "/", nil)

	b.RunParallel(func(p *testing.PB) {
		w := httptest.NewRecorder()
		for p.Next() {
			w.Body.Reset()
			indexHandler2(w, req)
		}
	})

	printMemStats(b)
}

func BenchmarkIndexHandler3(b *testing.B) {
	// 创建请求
	req, _ := http.NewRequest("GET", "/", nil)

	b.RunParallel(func(p *testing.PB) {
		w := httptest.NewRecorder()
		for p.Next() {
			indexHandler3(w, req)
		}
	})

	printMemStats(b)
}

func printMemStats(b *testing.B) {
	// 获取基准测试后的GC统计信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 计算GC次数
	gcCount := m.NumGC
	gcPauseTotalNs := m.PauseTotalNs
	alloc := m.Alloc
	totalAlloc := m.TotalAlloc

	b.Logf("test %s N: %d, biz size:%.2fKB, GC count: %v, pause: %.2fms, alloc: %.2fMB, total heap: %.2fMB",
		b.Name(),
		b.N,
		toKb(SIZE),
		gcCount,
		toMs(gcPauseTotalNs),
		toMb(alloc),
		toMb(totalAlloc))
}

func toKb(b uint64) float64 {
	return float64(b) / 1024.0
}

func toMb(b uint64) float64 {
	return float64(b) / 1024.0 / 1024.0
}

func toMs(b uint64) float64 {
	return float64(b) / 1000.0 / 1000.0
}
