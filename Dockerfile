# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.18-buster  AS build

WORKDIR /app

COPY go.* ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /server

##
## Deploy
##
FROM  gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=build /server .

COPY --from=busybox:1.35.0-uclibc /bin/sh /bin/sh

# create directory /var/pubsub/data/
RUN mkdir -p /var/pubsub/data/

EXPOSE 8080

ENTRYPOINT ["./server"]
