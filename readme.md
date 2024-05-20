# Set the GOPROXY environment variable
export GOPROXY=https://goproxy.io,direct
# Set environment variable allow bypassing the proxy for specified repos (optional)
export GOPRIVATE=git.mycompany.com,github.com/my/private

# 查看环境变量
go env

# 修改运行中容器的配额
docker container update 29cac6d4108f --cpus="4" --memory="8g" --memory-swap="-1"

# 查看容器运行状态
docker stats 29cac6d4108f