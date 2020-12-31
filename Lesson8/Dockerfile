FROM golang:1.14

WORKDIR /
#COPY . .
COPY ./api/api /bin
COPY ./process/process /bin


ENV CGO_ENABLED 0
ENV GO111MODULE on
RUN make server
#RUN ls -l
#RUN ls -l /api
#RUN CGO_ENABLED=0 go build -v -mod=vendor -ldflags "-s" -o /bin/api /api/.
#RUN CGO_ENABLED=0 go build -v -mod=vendor -ldflags "-s" -o /bin/process /process

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /bin ./bin