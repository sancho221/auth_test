FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache \
    curl \
    protoc \
    protobuf-dev

COPY go.mod go.sum ./
RUN go mod download

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY proto/ ./proto/
RUN mkdir -p pkg/pb && \
    protoc --go_out=./pkg/pb --go-grpc_out=./pkg/pb proto/auth.proto

COPY . .

RUN go build -o /auth-service ./cmd/main.go

CMD ["/auth-service"]