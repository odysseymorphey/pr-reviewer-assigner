FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o server ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY ./config ./config

EXPOSE 8080

ENTRYPOINT [ "./server" ]