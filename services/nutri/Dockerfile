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

#
#FROM scratch
#
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
#
WORKDIR /app
#COPY --from=builder /app/nutri /app/nutri

CMD ["./nutri"]

#
#FROM golang:1.21.5-alpine
## Update and install necessary packages
#RUN apk update && apk add git
#RUN apk add --no-progress --no-cache gcc musl-dev
## Set environment variables
#ENV CGO_ENABLED=1
#ENV GO111MODULE=on
#ENV GOPATH /go
#ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
## Create and set working directory
#RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
#WORKDIR $GOPATH/src/devops
## Copy the source code
#COPY . .
## Download dependencies
#RUN go mod download
## Build the application
#RUN GOOS=linux go build -tags musl -ldflags '-extldflags "-static"' -o nutri nutri/cmd/nutri.go
## Set the entry point for the container
#ENTRYPOINT ["./nutri"]
