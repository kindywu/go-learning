# Set the GOPROXY environment variable

export GOPROXY=https://goproxy.io,direct

# Set environment variable allow bypassing the proxy for specified repos (optional)

export GOPRIVATE=git.mycompany.com,github.com/my/private

# Set environment variable in windows 11

$env:GOPROXY = "https://goproxy.io"

$env:GOPROXY = "https://goproxy.cn"

# 查看环境变量

go env

# 修改运行中容器的配额

docker container update 29cac6d4108f --cpus="4" --memory="8g" --memory-swap="-1"

# 查看容器运行状态

docker stats 29cac6d4108f

# 开发时自动热启动

go install github.com/cosmtrek/air@latest

# cobra

- go get -u github.com/spf13/cobra/cobra
- go install github.com/spf13/cobra-cli@latest

![alt text](image.png)

- 要将微秒（μs）转换为纳秒（ns），只需将 μs 乘以 1000，因为 1 微秒等于 1000 纳秒。因此，50-100 μs 的范围等于 50,000-100,000 纳秒。

# 延迟（Latency）和带宽（Bandwidth）都是衡量存储系统或通信系统性能的重要指标，但它们代表的是不同方面的速度。

- 延迟（Latency）：
  延迟指的是从发出请求到接收到响应所花费的时间。在存储系统中，延迟表示存储设备或存储系统对数据请求的响应速度。更低的延迟意味着系统能够更快地响应数据请求，因此延迟越低越好。
  例如，在存储系统中，延迟通常表示从发出读取请求到数据可用的时间，或者从发出写入请求到写入完成的时间。
- 带宽（Bandwidth）：
  带宽表示在单位时间内能够传输的数据量。在存储系统中，带宽表示存储系统能够传输的数据速率。更高的带宽意味着系统能够在单位时间内传输更多的数据，因此带宽越高越好。
  例如，在存储系统中，带宽通常表示在一秒内能够从存储设备读取或写入的数据量。

# 配置 GOPROXY 环境变量
export GOPROXY=https://goproxy.io,direct

秒 (second) 1
毫秒（Millisecond）1000
微秒（Microsecond）1000 * 1000
纳秒（Nanosecond） 1000 * 1000 * 1000