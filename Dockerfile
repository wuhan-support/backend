FROM golang:1.13-alpine

MAINTAINER homilly

WORKDIR $GOPATH/src/

ADD * $GOPATH/src/

COPY config.example.yml config.yml

CMD ["/bin/bash","go","run","."]
