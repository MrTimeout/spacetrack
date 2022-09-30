# golang:alpine with version 3.16 (last alpine version also)
ARG USER=basket1
ARG UID=1001
FROM golang@sha256:c9a90742f5457fae80d8f9f1c9fc6acd6884c749dc6c5b11c22976973564dd4f as base

ARG USER
ARG UID

RUN adduser -h /home/${USER} -D -u ${UID} -s /sbin/nologin -g "" ${USER} && \
  mkdir --parent /go/src/github.com/MrTimeout/go-spacetrack && \
  apk add --no-cache --update ca-certificates tzdata git && update-ca-certificates

WORKDIR /go/src/github.com/MrTimeout/spacetrack

COPY . .

RUN go mod download && go mod verify && \
  CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-w -s" -o /go/bin ./... && \
  chown -R "${UID}" /go

USER "${USER}"

CMD [\
  "/go/bin/go-spacetrack", \
  "--work-dir", "/tmp/upload/basket1/products/automatic", \
  "--format", "xml", \
  "--rest-call", "all", \
]