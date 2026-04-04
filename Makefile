.PHONY: build run clean migrate docs

# 构建
build:
	go build -o bin/server cmd/server/main.go cmd/server/app.go

# 运行
run:
	go run cmd/server/main.go cmd/server/app.go

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

# 生成 Swagger 文档
docs:
	swag init -g main.go -d cmd/server,internal/handler,internal/pkg/req,internal/pkg/response -o cmd/server/docs
