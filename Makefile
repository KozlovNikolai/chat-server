LOCAL_BIN:=$(CURDIR)/bin

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	GOBIN=$(LOCAL_BIN) $(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-chat-api

generate-chat-api:
	mkdir -p pkg/chat_v1
	protoc --proto_path api/chat_v1 \
	--go_out=pkg/chat_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/chat_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/chat_v1/chat.proto

build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc-server/main.go

copy-to-server:
	scp service_linux root@5.159.101.123:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/msc/chat-server:v0.0.1 .
	docker login -u token -p CRgAAAAAD6l6BhyNZkQouCVYue4xorwBot6D5eZ6 cr.selcloud.ru/msc
	docker push cr.selcloud.ru/msc/chat-server:v0.0.1

run-into-server:
	docker pull cr.selcloud.ru/msc/chat-server:v0.0.1
	docker run -p 50051:50051 cr.selcloud.ru/msc/chat-server:v0.0.1

docker-local-build-and-run:
	-docker stop chat
	-docker rm chat
	-docker rmi chat:v0.0.1
	docker buildx build --no-cache --platform linux/amd64 -t chat:v0.0.1 .
	docker run --name chat -p 50051:50051 chat:v0.0.1
