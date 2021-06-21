# Builder

FROM golang:1.16 AS builder

WORKDIR /data
COPY . .

RUN CGO_ENABLED=0 go build -v

# Runtime

FROM scratch

COPY --from=builder /data/grpc-json-proxy /

EXPOSE 7001

CMD ["/grpc-json-proxy"]