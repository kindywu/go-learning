package main

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
)

func BenchmarkIndexHandler(b *testing.B) {
	// 获取初始的GC统计信息
	var m0, m1 runtime.MemStats
	runtime.ReadMemStats(&m0)

	b.RunParallel(func(p *testing.PB) {
		// 创建一个请求
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			b.Fatal("Error creating request: ", err)
		}

		// 创建响应记录器
		w := httptest.NewRecorder()

		for p.Next() {
			indexHandler(w, req)
		}
	})

	// 获取基准测试后的GC统计信息
	runtime.ReadMemStats(&m1)

	// 计算GC次数
	gcCount := m1.NumGC - m0.NumGC
	b.Logf("GC count: %v", gcCount)
}
