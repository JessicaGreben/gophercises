# syntax=docker/dockerfile:1

## Build
FROM golang:1.18-buster AS build

WORKDIR /app

COPY . ./

RUN go build -o /exercise1

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /exercise1 /exercise1
COPY --from=build /app/quizzes /quizzes

USER nonroot:nonroot

ENTRYPOINT ["/exercise1"]
