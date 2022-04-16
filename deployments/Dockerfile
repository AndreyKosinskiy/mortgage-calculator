# syntax=docker/dockerfile:1

# Build
FROM golang:1.16-alpine AS build
WORKDIR /app


COPY /go.mod ./
COPY /go.sum ./

RUN go mod download

COPY ../ ./

RUN cd cmd/ && GOARCH=amd64 CGO_ENABLED=0 GOOS=linux go build -v -o ../build/backend

#Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build app/build/backend /

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/backend"]