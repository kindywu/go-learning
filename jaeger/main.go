package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	// 创建一个新的上下文。
	// 创建一个 OTLP HTTP 导出器。
	// 这里，我们将数据发送到运行在本地的 Jaeger 实例，默认端口为 4318。
	// 确保在程序结束时关闭导出器
	// 设置全局的 TracerProvider。
	// 创建一个新的 Span。
	// 在 Span 中执行一些工作。
	// 模拟工作负载
	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL("http://jaeger-jaeger-1:4318"))
	if err != nil {
		log.Fatalf("failed to create Jaeger exporter: %v", err)
	}
	defer exporter.Shutdown(ctx)
	println("连接成功")
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("local-service"), // 设置服务名称
	)

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	println("创建Tracer成功")

	go callService(ctx)

	// Respect OS stop signals.
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-c
}

func callService(ctx context.Context) {
	tracer := otel.Tracer("example-tracer")
	parentCtx, span := tracer.Start(ctx, "call-main-service")
	defer span.End()

	span.SetAttributes(attribute.String("user.name", "kindywu"))
	log.Printf("ctx:%+v\n", span.SpanContext())

	span.AddEvent("call-main-service-start")
	time.Sleep(1 * time.Second)

	println("创建Span成功")
	span.AddEvent("call-main-service-end")

	time.Sleep(1 * time.Second)
	println("调用call-main-service成功")

	go callOtherService(1, parentCtx)

}

func callOtherService(id int, ctx context.Context) {

	serviceName := fmt.Sprintf("call-other-service %d", id)
	// log.Printf("ctx:%+v\n", ctx)
	tracer := otel.Tracer("example-tracer")
	childCtx, span := tracer.Start(ctx, serviceName)
	log.Printf("ctx:%+v\n", span.SpanContext())
	defer span.End()
	println("创建ChildSpan成功")
	span.AddEvent(serviceName + "-start")
	time.Sleep(3 * time.Second)
	span.AddEvent(serviceName + "-end")
	println("调用" + serviceName + "成功")

	go callHttpService(childCtx)
}

func callHttpService(ctx context.Context) {

	// 创建一个新的HTTP请求。
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	// req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/users/123", nil)
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:3000/users/1234", nil)

	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	propagator := otel.GetTextMapPropagator()
	// 使用传播器将追踪上下文注入到HTTP请求头部。
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	fmt.Printf("Inject %v %v \n", ctx, propagator.Fields())

	// 发起HTTP请求到另一个服务。
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// 这里可以处理响应，例如检查状态码或读取响应体。
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK status: %v", resp.Status)
	}

	// 示例：读取响应体内容。
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
	}
	log.Printf("Response body: %s", bodyBytes)
}
