FROM golang:alpine AS builder

LABEL stage=gobuilder

RUN apk update
RUN apk add \
    gcc \
    musl-dev \
    # explicitly install SASL package
    cyrus-sasl-dev

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

ENV CGO_ENABLED 1
ENV GO111MODULE=on
ENV CGO_LDFLAGS="-lsasl2"

COPY . .

RUN go build -tags musl -ldflags="-s -w" -o /app/nutri cmd/nutri.go

WORKDIR /app

CMD ["./nutri"]

