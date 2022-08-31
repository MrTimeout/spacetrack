# golang:alpine with version 3.16 (last alpine version also)
FROM golang@sha256:c9a90742f5457fae80d8f9f1c9fc6acd6884c749dc6c5b11c22976973564dd4f as base

ENV USER=spacetrack-client \
  UID=10001

RUN adduser -h /home/${USER} -D -u ${UID} -s /sbin/nologin -g "" ${USER} && \
  mkdir --parent /go/src/github.com/MrTimeout/spacetrack && \
  apk add --no-cache --update ca-certificates tzdata git && update-ca-certificates

WORKDIR /go/src/github.com/MrTimeout/spacetrack

COPY . .

RUN go mod download && go mod verify && \
  CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-w -s" -o /go/bin ./...

FROM alpine@sha256:686d8c9dfa6f3ccfc8230bc3178d23f84eeaf7e457f36f271ab1acc53015037c as pre

COPY --from=base /usr/share/zoneinfo/ /usr/share/zoneinfo/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=base /home/spacetrack-client /home/spacetrack-client

COPY --from=base /go/bin/spacetrack /go/bin/spacetrack

USER spacetrack-client:spacetrack-client

CMD [\
  "/go/bin/spacetrack", "gp", "--log-level", "debug", \
  "--log-file", "/tmp/spacetrack.json", \
  "--config", "/home/spacetrack-client/.spacetrack.yaml", \
  "--work-dir", "/tmp", \
  "--limit", "5", \
  "--skip", "2", \
  "--format", "json", \
  "--filter", "decay_date<>null-val", \
  "--filter", "epoch<now-30", \
  "--orderby", "norad_cat_id", \
  "--sort", "asc"\
  ]
