# 默认目标
.PHONY: build test install clean lint fmt vet

# 构建二进制文件
build:
	go mod download
	go build -o saturn_cli ./examples/client/client.go
	go build -o saturn_svr ./examples/server/server.go
	chmod +x .githooks/*
	git config core.hooksPath .githooks

# 运行测试
test:
	go test -count 1 -v ./... -gcflags "all=-N -l" && go test -race -v ./... -gcflags "-l"

# 安装二进制文件到系统路径
install:
	go install ./examples/client/client.go
	go install ./examples/server/server.go

# 清理构建产物
clean:
	rm -f saturn_cli saturn_svr

# 代码格式化
fmt:
	go fmt ./...

# 代码检查
vet:
	go vet ./...

# 使用golangci-lint进行代码检查
lint:
	golangci-lint run --no-config

# 构建发布版本（带优化标记）
release:
	go build -ldflags "-s -w" -o saturn_cli ./examples/client/client.go
	go build -ldflags "-s -w" -o saturn_svr ./examples/server/server.go

# 运行基准测试
bench:
	go test -bench=. -benchmem ./...