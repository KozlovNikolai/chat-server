#FROM golang:1.20.3-alpine AS builder
FROM golang:1.21.7-alpine AS builder

COPY . /github.com/KozlovNikolai/chat-server/source/
WORKDIR /github.com/KozlovNikolai/chat-server/source/

RUN go mod download
RUN go build -o ./bin/chat-server cmd/grpc-server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/KozlovNikolai/chat-server/source/bin/chat-server .
COPY --from=builder /github.com/KozlovNikolai/chat-server/source/*.env .


CMD ["./chat-server -config-path local.env"]