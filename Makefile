.PHONY: build run clean migrate

# 构建
build:
	go build -o bin/server cmd/server/main.go

# 运行
run:
	go run cmd/server/main.go

# 清理
clean:
	rm -rf bin/

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
vet:
	go vet ./...

# 下载依赖
tidy:
	go mod tidy
