# Builder

FROM golang AS builder

WORKDIR /data
COPY . /data

RUN go mod download \
	&& go build

# Runtime

FROM debian:stretch-slim

MAINTAINER vincent <vincent.h.cui@gmail.com>

COPY --from=builder /data/grpc-json-proxy /

EXPOSE 7001

CMD ["/grpc-json-proxy"]