FROM golang:1.15-alpine

RUN apk --update add build-base

WORKDIR /go/src/github.com/wamuir/go-jsonapi-server
COPY . .

RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-linkmode external -extldflags "-static"' -o /go/bin/go-jsonapi-server

USER 65534:65534

ENTRYPOINT ["/go/bin/go-jsonapi-server"]

EXPOSE 8080/tcp
