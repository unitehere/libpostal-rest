FROM golang:1.11-alpine AS builder

WORKDIR /go/src/app
COPY . /go/src/app
RUN apk --update --no-cache add --virtual .build-deps autoconf automake curl \
      gcc g++ libtool make pkgconfig git \
      && sh ./install-libpostal.sh \
      && go get ./... \
      && go build -o main

FROM alpine

COPY --from=builder /usr/local/lib/libpostal.so.1.0.0 /usr/local/lib/libpostal.so.1.0.0
RUN cd /usr/local/lib \
    && ln -s libpostal.so.1.0.0 libpostal.so.1 \
    && ln -s libpostal.so.1.0.0 libpostal.so
COPY --from=builder /usr/local/lib/libpostal.la /usr/local/lib/libpostal.la
COPY --from=builder /usr/local/lib/libpostal.a /usr/local/lib/libpostal.a

WORKDIR /app
COPY --from=builder /go/src/app/main .

EXPOSE 8080
ENTRYPOINT /app/main

