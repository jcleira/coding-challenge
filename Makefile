test:
	go test -v --race ./...

run:
	go run main.go

lint:
	golangci-lint run

mocks:
	mockery --all --case snake --disable-version-string

.PHONY: mocks
