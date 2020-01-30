FROM golang
MAINTAINER  homilly

COPY config.example.yml config.yml

ADD src /go/src/backend

WORKDIR /go/src/backend

EXPOSE 80

CMD ["/bin/bash","go","run","."]
