FROM golang:1.15 AS builder
RUN apt update ; apt install upx-ucl -y ; apt clean
WORKDIR /go/src/github.com/kronostechnologies/slack-hook/
COPY * ./
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o slack-hook . && upx --best slack-hook

FROM scratch
COPY --from=builder /go/src/github.com/kronostechnologies/slack-hook/slack-hook /bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/bin/slack-hook"]
