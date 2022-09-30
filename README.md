# Go Spacetrack

## Build image

This would be a standard build, with all the defaults (USER AND UID)

`docker image build --file ./Dockerfile --tag go-spacetrack:0.3 --no-cache .`

We can modify the _USER_ and _UID_ like we want it.

`docker image build --file ./Dockerfile --tag go-spacetrack:0.3 --no-cache --build-arg USER=MrTimeout --build-arg UID=1002 .`

We can also append our repository and username in that repository. In my case, I'm going to upload the image to the Docker hub under the username _estenoesmiputonombre_

`docker image build --file ./Dockerfile --tag estenoesmiputonombre/spacetrack:0.3 --no-cache .`

## Run image

`${USER}` must be the user that you have passed to `--build-arg USER` or `basket1` by default

```sh
# Remove the volume if you want a fresh start
> docker volume rm spacetrack-data

# Execute this script if you want to set the volume from a fresh start
> docker container run --name ccc --rm -it -e "USER=basket1" -e "UID=1001" --volume spacetrack-data:/tmp/upload \
  alpine:latest /bin/sh -c "adduser -h /home/\${USER} -D -u \${UID} -s /sbin/nologin -g \"\" \${USER} && chown -R \${USER}:\${USER} /tmp/upload"

> docker container run --name spacetrack-client --detach \
  --volume ${PWD}/spacetrack.yml:/home/${USER}/spacetrack.yml \
  --volume spacetrack-data:/tmp/upload \
  estenoesmiputonombre/spacetrack:0.3
```

We can also modify son options at runtime

```sh
> docker container run --name spacetrack-client --detach \
  --volume ${PWD}/spacetrack.yml:/home/${USER}/spacetrack.yml \
  --volume spacetrack-data:/tmp/spacetrack \
  estenoesmiputonombre/spacetrack:0.3 /go/bin/go-spacetrack --format=xml --rest-call=tle --work-dir=/tmp/spacetrack
```
