FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /auth-service ./cmd/main.go

EXPOSE 8082

CMD [ "/auth-service" ]