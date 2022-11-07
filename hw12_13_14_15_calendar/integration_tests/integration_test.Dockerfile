FROM golang:1.16-alpine

RUN apk add --no-cache make curl gcc libc-dev

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ENV CGO_ENABLED=0

CMD go test -v -tags integration ./integration_tests