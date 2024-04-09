FROM golang:1.22-alpine3.18 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build cmd/main.go

# Path: Dockerfile
FROM alpine:3.18 as runner

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
