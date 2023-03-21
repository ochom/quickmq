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

# Busybox for scripts
FROM busybox:1.35.0-uclibc as busybox

##
## Deploy
##
FROM  gcr.io/distroless/base-debian11

WORKDIR /app

COPY --from=build /server .

# copy busybox shell to /bin/sh
COPY --from=busybox /bin/sh /bin/sh

# copy busybox mkdir to /bin/mkdir

COPY --from=busybox /bin/mkdir /bin/mkdir

# create directory /var/pubsub/data/
RUN mkdir -p /var/pubsub/data/

EXPOSE 8080

ENTRYPOINT ["./server"]
