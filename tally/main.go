package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
)

func main() {
	// 创建一个新的 tally Scope
	reporter := prometheus.NewReporter(prometheus.Options{
		OnRegisterError: func(err error) {
			log.Fatalln("Error registering Prometheus metric", err)
		},
	})
	prometheusScope, prometheusCloser := tally.NewRootScope(tally.ScopeOptions{
		Prefix:         "demo",
		Tags:           map[string]string{"node_name": "demo"},
		CachedReporter: reporter,
		Separator:      "_",
		SanitizeOptions: &tally.SanitizeOptions{
			NameCharacters: tally.ValidCharacters{
				Ranges:     tally.AlphanumericRange,
				Characters: tally.UnderscoreCharacters,
			},
			KeyCharacters: tally.ValidCharacters{
				Ranges:     tally.AlphanumericRange,
				Characters: tally.UnderscoreCharacters,
			},
			ValueCharacters: tally.ValidCharacters{
				Ranges:     tally.AlphanumericRange,
				Characters: tally.UnderscoreCharacters,
			},
			ReplacementCharacter: tally.DefaultReplacementCharacter,
		},
	}, time.Duration(1*time.Second))
	defer prometheusCloser.Close()

	CORSHeaders := handlers.AllowedHeaders([]string{"Content-Type", "User-Agent"})
	CORSOrigins := handlers.AllowedOrigins([]string{"*"})
	CORSMethods := handlers.AllowedMethods([]string{"GET", "HEAD"})
	// 注册 Prometheus 的 HTTP 处理程序
	handler := handlers.CORS(CORSHeaders, CORSOrigins, CORSMethods)(reporter.HTTPHandler())
	http.Handle("/metrics", handler)

	// 启动 HTTP 服务器来暴露 metrics
	go func() {
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for t := range ticker.C {
			value := rand.Float64() * 1000
			fmt.Println("Tick at", t, value)
			prometheusScope.Gauge("lua_runtimes").Update(value)
		}
	}()

	// Respect OS stop signals.
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-c
}
