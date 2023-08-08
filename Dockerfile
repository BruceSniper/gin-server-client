# FROM ... AS builder : 表示依赖的镜像只是使用在编译阶段
FROM golang:1.18.1 AS builder

# 编译阶段的工作目录，也可以作为全局工作目录
WORKDIR /app

# 把当前目录的所有内容copy到 WORKDIR指定的目录中
COPY . .

# 定义go build的工作环境，
# 例如GOOS=linux、GOARCH=amd64，
# 这样编译出来的 'main可执行文件' 就只能在linux的amd64架构中使用
ARG TARGETOS
ARG TARGETARCH

# 执行go build； --mount：在执行build时，会把/go 和 /root/.cache/go-build 临时挂在到容器中
RUN --mount=type=cache,target=/go --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main

FROM alpine:3.14.0

# 把执行builder阶段的结果 /app/main拷贝到/app中
COPY --from=builder /app/main /app
EXPOSE 8080
# 运行main命令，启动项目
# /app/main 指向RUN命令的 go build -o main的结果
ENTRYPOINT ["/app/main"]