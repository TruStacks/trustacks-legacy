.PHONY: test build clean

test:
	@go test ./... -v -race
	@golangci-lint run 

build:
	@CGO_ENABLED=0 go build -o tsctl -ldflags "-X main.cliVersion=0.2.5" ./cmd 

clean:
	@rm tsctl