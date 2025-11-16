FROM golang:1.25 as builder

RUN go install go.k6.io/xk6/cmd/xk6@latest

RUN /go/bin/xk6 build \
    --with github.com/grafana/xk6-dashboard@latest \
    --output /k6

FROM alpine:latest

RUN apk add --no-cache ca-certificates

COPY --from=builder /k6 /usr/bin/k6

WORKDIR /scripts

ENTRYPOINT ["k6"]
