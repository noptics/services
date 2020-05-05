FROM golang:1.14 as Builder

RUN mkdir -p /go/src/github.com/noptics/services
ADD . /go/src/github.com/noptics/services

WORKDIR /go/src/github.com/noptics/services/cmd/nats-proxy

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o nats-proxy

FROM alpine:3.9

RUN apk add --no-cache curl bash ca-certificates

COPY --from=builder /go/src/github.com/noptics/services/cmd/nats-proxy/nats-proxy /nats-proxy

CMD ["/nats-proxy"]