build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./go-spacetrack ./...

debug:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./go-spacetrack -gcflags="all=-N -l" ./...
	@dlv --listen=:2345 --headless=true --api-version=2 exec ./go-spacetrack -- --config-file /tmp/spacetrack.yml

docker-build:
	@docker image build --file ./Dockerfile --tag estenoesmiputonombre/spacetrack:0.3 --no-cache --build-arg USER=basket1 --build-arg UID=1001 .

docker-run:
	@docker container run --name spacetrack --detach --volume ${PWD}/spacetrack.yml:/home/basket1/spacetrack.yml --volume spacetrack-data:/tmp/upload estenoesmiputonombre/spacetrack:0.3