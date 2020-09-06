FROM golang:alpine AS builder
ARG VERSION
WORKDIR /go/src/sheet-CRUD
COPY . .
WORKDIR /go/src/sheet-CRUD/main
RUN apk add --no-cache --virtual .build-deps git libc6-compat build-base; \
    go get github.com/go-delve/delve/cmd/dlv
RUN go mod download
RUN go get -d -v; \
    go install -v; \
    go build -gcflags "all=-N -l" -ldflags "-X main.version=${VERSION}" -o app;

FROM alpine:latest
WORKDIR /root
COPY --from=builder /go/src/sheet-CRUD/main/app .

CMD ["./app"]
