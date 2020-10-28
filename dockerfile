FROM golang:alpine AS builder
ARG VERSION
RUN apk add --no-cache --virtual .build-deps git libc6-compat build-base
WORKDIR /sheet-CRUD

COPY . .
RUN go mod download
WORKDIR /sheet-CRUD/main
RUN go build -gcflags "all=-N -l" -ldflags "-X main.version=${VERSION}" -o app;

FROM alpine:latest
WORKDIR /root
COPY --from=builder /sheet-CRUD/main/app .

COPY ./swagger-ui /root/swagger-ui

CMD ["./app"]
