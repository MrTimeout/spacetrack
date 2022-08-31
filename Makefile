build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bin -ldflags="-w -s" ./...

build-local:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ${GOPATH}/bin -ldflags="-w -s" ./...

clean:
	@rm ./bin/spacetrack*

docker_tpz:
	@docker image build --file ./Dockerfile.tpz --tag spacetrack:tpz --no-cache .
	@docker container rm -f spacetrack-tpz
	@docker container run --name spacetrack-tpz --detach -v spacetrack:/tmp/upload -v ${PWD}/spacetrack.yaml:/root/.spacetrack.yaml spacetrack:tpz

install:
	GOARCH=amd64 GOOS=linux go install ./...