TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOPATH?=$$(go env GOPATH)
GRPC_GATEWAY_PATH?=$$(go list -m -u -f '{{ .Dir }}' all | grep 'github.com/grpc-ecosystem/grpc-gateway@')

default: test

lint:
	golangci-lint run ./...

fmt:
	gofmt -w $(GOFMT_FILES)

cover:
	go test $(TEST) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out

test:
	@sh -c "go test ./... -race -timeout=2m -parallel=4"

proto:
	protoc --proto_path=grpc/protobuf/:$(GRPC_GATEWAY_PATH)/third_party/googleapis/ \
		--go_out=plugins=grpc:. \
		--grpc-gateway_out=logtostderr=true:. \
		health.proto

.PHONY: fmt lint test
