package main

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	fiberV2 "github.com/gofiber/fiber/v2"
	fiberV3 "github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

// go test -bench=BenchmarkIndexHandler1 -benchmem -benchtime=10s -v
// go test -bench=BenchmarkIndexHandler2 -benchmem -benchtime=10s -v
// go test -bench=BenchmarkIndexHandler3 -benchmem -benchtime=10s -v
// go test -bench=BenchmarkIndexHandler4 -benchmem -benchtime=10s -v
// go test -bench=BenchmarkFiberHandlerV2 -benchmem -benchtime=10s -v
// go test -bench=BenchmarkFiberHandlerV3 -benchmem -benchtime=10s -v

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
	req, _ := http.NewRequest("GET", "/1", nil)
	w := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		indexHandler1(w, req)
	}
	printMemStats(b)
}

func BenchmarkIndexHandler2(b *testing.B) {
	req, _ := http.NewRequest("GET", "/2", nil)
	w := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		indexHandler2(w, req)
	}
	printMemStats(b)
}

func BenchmarkIndexHandler3(b *testing.B) {
	req, _ := http.NewRequest("GET", "/3", nil)
	w := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		indexHandler3(w, req)
	}
	printMemStats(b)
}

func BenchmarkIndexHandler4(b *testing.B) {
	req, _ := http.NewRequest("GET", "/4", nil)
	w := httptest.NewRecorder()
	for n := 0; n < b.N; n++ {
		indexHandler4(w, req)
	}
	printMemStats(b)
}

func BenchmarkFiberHandlerV2(b *testing.B) {
	app := fiberV2.New()

	app.Get("/", func(c *fiberV2.Ctx) error {
		return c.SendString("Hello, World!")
		// return c.SendStatus(fiber.StatusOK)
	})

	h := app.Handler()
	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(fiberV2.MethodGet)
	fctx.Request.SetRequestURI("/")

	for n := 0; n < b.N; n++ {
		h(fctx)
	}
	printMemStats(b)
}

// var msg2 = Msg{Message: "Hello, World!"}

func BenchmarkFiberHandlerV3(b *testing.B) {
	app := fiberV3.New()

	app.Get("/", func(c fiberV3.Ctx) error {
		// return c.JSON(fiberV3.Map{"message": "Hello, World!"})
		// return c.JSON(msg2)
		return c.JSON(Msg{Message: "Hello, World!"})
		// return c.SendString("Hello, World!")
		// return c.SendStatus(fiber.StatusOK)
	})

	h := app.Handler()
	fctx := &fasthttp.RequestCtx{}
	fctx.Request.Header.SetMethod(fiberV3.MethodGet)
	fctx.Request.SetRequestURI("/")

	for n := 0; n < b.N; n++ {
		h(fctx)
	}
	printMemStats(b)
}

func printMemStats(b *testing.B) {
	// 获取基准测试后的GC统计信息
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// 计算GC次数
	numGC := m.NumGC
	pauseTotalNs := m.PauseTotalNs
	alloc := m.Alloc
	totalAlloc := m.TotalAlloc

	b.Logf("test N: %d, biz_size:%.2fKB, GC: %v, total_pause: %.2fms, alloc: %.2fMB, total_alloc: %.2fMB",
		b.N,
		toKb(SIZE),
		numGC,
		toMs(pauseTotalNs),
		toMb(alloc),
		toMb(totalAlloc))

	once.Do(func() {
		printMaxGoroutine(b)
	})
}

var max_numGoroutines atomic.Int32

func printMaxGoroutine(b *testing.B) {
	b.Log("start printMaxGoroutine")
	go func() {
		for {
			numGoroutine := runtime.NumGoroutine()
			if int(max_numGoroutines.Load()) < numGoroutine {
				max_numGoroutines.Store(int32(numGoroutine))
				b.Logf("max num goroutines: %d", max_numGoroutines.Load())
				time.Sleep(1 * time.Millisecond)
			}

		}
	}()
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
