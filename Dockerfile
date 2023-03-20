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
FROM golang:1.18-buster AS deploy

WORKDIR /app

COPY --from=build /server .

RUN mkdir -p /var/pubsub/data/

EXPOSE 8080

ENTRYPOINT ["./server"]
