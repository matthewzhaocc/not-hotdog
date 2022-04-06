FROM golang:latest AS builder
RUN mkdir /build
WORKDIR /build
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .

FROM alpine:3.15.4
RUN apk update && apk add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /build/not-hotdog not-hotdog
COPY --from=builder /build/index.html index.html
CMD ./not-hotdog
