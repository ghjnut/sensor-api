FROM golang:1.20-alpine3.17 as builder

#RUN apk add git

WORKDIR /go/src/github.com/ghjnut/sensor-api

COPY go.mod ./
# TODO add back when we have deps
#COPY go.sum ./

RUN go mod download

COPY . .

RUN go build

FROM alpine:3.17

COPY --from=builder /go/bin/* /usr/local/bin/

#VOLUME /etc/pingwave.hcl

ENTRYPOINT ["sensor-api"]
