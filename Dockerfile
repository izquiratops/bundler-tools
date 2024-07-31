# Use this image to build the executable
FROM golang:1.18-alpine AS build

WORKDIR /go/src/github.com/izquiratops/bundler-tools
COPY . /go/src/github.com/izquiratops/bundler-tools/

RUN apk add --no-cache git ca-certificates make bash
RUN /usr/bin/env bash -c make install

# Final image containing the executable from the previous step
FROM alpine:3

COPY --from=build /go/bin/bundler-tools /usr/bin/bundler-tools
COPY "container/entrypoint.sh" "/init.sh"

ENTRYPOINT ["/init.sh"]