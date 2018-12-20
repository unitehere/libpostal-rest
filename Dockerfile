FROM golang:1.11-alpine

WORKDIR /go/src/app
COPY . /go/src/app
RUN apk --update --no-cache add --virtual .build-deps autoconf automake curl \
      gcc g++ libtool make pkgconfig git \
      && sh ./install-libpostal.sh \
      && go get ./... \
      && go build -o main

EXPOSE 8080

ENTRYPOINT /go/src/app/main

