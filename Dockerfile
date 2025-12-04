FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./src/main.go 

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main /app/main
COPY config ./config
COPY templates ./templates
COPY migration ./migration

RUN chmod +x /app/main

CMD ["./main"]