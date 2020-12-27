FROM golang:1.14

WORKDIR /go/src/github.com/alxmsl/geekbrains.queues/
COPY . .

ENV CGO_ENABLED 0
ENV GO111MODULE on
RUN make server



FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/alxmsl/geekbrains.queues/bin ./bin