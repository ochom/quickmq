# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.22-alpine  AS build

WORKDIR /app

COPY go.* ./
COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o /server ./cmd/main.go

##
## Build UI
##

FROM oven/bun:1.0 AS ui

WORKDIR /app

COPY web ./

RUN bun install

RUN bun run build

##
## Deploy
##
FROM busybox:1.35.0-uclibc AS deploy 

WORKDIR /

RUN mkdir -p /web

COPY --from=build /server .
COPY --from=ui /app/build /web/build

EXPOSE 16321
EXPOSE 6321

ENTRYPOINT ["/server"]
