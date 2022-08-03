build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bin -ldflags="-w -s" ./...

build-local:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ${GOPATH}/bin -ldflags="-w -s" ./...

clean:
	@rm ./bin/spacetrack*

install:
	GOARCH=amd64 GOOS=linux go install ./...