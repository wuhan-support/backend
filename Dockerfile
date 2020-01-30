FROM golang:1.13-alpine

MAINTAINER homilly

WORKDIR $GOPATH/src/

ADD * $GOPATH/src/

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk add git

COPY config.example.yml config.yml

ENTRYPOINT go run .
