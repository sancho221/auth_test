FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache curl

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /auth-service ./cmd/main.go

CMD ["/auth-service"]