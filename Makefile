build:
	go mod download
	go build  -o saturn_cli ./examples/client/client.go
	go build  -o saturn_svr ./examples/server/server.go
	chmod +x .githooks/*
	git config core.hooksPath .githooks

test:
	go test -count 1 -v ./... -gcflags "all=-N -l" && go test -race -v ./... -gcflags "-l"

install:

