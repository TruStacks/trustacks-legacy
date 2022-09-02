.PHONY: test build clean

test:
	@go test ./... -v -race
	@golangci-lint run 

test_firebase:
	@firebase emulators:exec --only firestore --project fake-project-id 'go test ./... -v -race -tags=test_firebase'

build:
	@CGO_ENABLED=0 go build -o tsctl ./cmd

clean:
	@rm tsctl