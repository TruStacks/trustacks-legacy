.PHONY: test build clean

test: export FIRESTORE_EMULATOR_HOST=localhost:54321
test:
	@go test ./... -v -race
	@golangci-lint run 

build:
	@CGO_ENABLED=0 go build -o tsctl ./cmd

clean:
	@rm tsctl