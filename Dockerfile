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
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /server .

EXPOSE 8080

ENTRYPOINT ["/server"]
