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

RUN go build -o /server .


##
## Deploy
##
FROM gcr.io/distroless/base-debian12:nonroot

WORKDIR /

COPY --from=build /server .

EXPOSE 16321
EXPOSE 6321

ENTRYPOINT ["/server"]
