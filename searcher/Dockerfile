FROM golang:1.22-alpine3.18 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN  apk add --no-cache gcc musl-dev

RUN go build cmd/main.go

# Path: Dockerfile
FROM alpine:3.18 as runner

WORKDIR /app

COPY --from=builder /app/main .

ENV PORT 8080

EXPOSE 8080

CMD ["./main"]
