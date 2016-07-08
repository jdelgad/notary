all: lint test
	go build

lint:
	gometalinter ./... --enable=test,gofmt,goimports,lll,misspell

test:
	go test --race ./...
