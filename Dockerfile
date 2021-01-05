FROM golang:latest AS builder
RUN mkdir /build
WORKDIR /build
COPY . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .

FROM busybox
RUN mkdir /app
WORKDIR /app
COPY --from=0 /build/not-hotdog not-hotdog
COPY --from=0 /build/index.html index.html
CMD ./not-hotdog
