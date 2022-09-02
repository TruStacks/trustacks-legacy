.PHONY: test build clean

test:
	@go test ./... -v -race
	@golangci-lint run 

test_firebase:
	@firebase emulators:exec --only firestore --project fake-project-id 'go test ./... -v -race -tags=test_firebase'

build:
	@CGO_ENABLED=0 go build -o tsctl -ldflags "-X main.cliVersion=0.2.5" ./cmd 

clean:
	@rm tsctl