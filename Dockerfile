FROM golang:latest AS builder
RUN mkdir /build
WORKDIR /build
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .

FROM alpine:latest AS certs
RUN apk update && apk add ca-certificates

FROM busybox
RUN mkdir /app
WORKDIR /app
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
COPY --from=builder /build/not-hotdog not-hotdog
COPY --from=builder /build/index.html index.html
CMD ./not-hotdog
