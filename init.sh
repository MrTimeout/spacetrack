#!/bin/bash
#
# This script is used to make sure that we have the proper environment to execute the spacetrack script and the plugins.
#
# We can customize this variables with the right values.
INGESTION_SPACETRACK_TAG="0.3"
INGESTION_SPACETRACK_VOLUME="ingestion_spacetrack"
INGESTION_SPACETRACK_FILE="/home/basket1/spacetrack.yml"
# Path to the file where our spacetrack configuration is.
INGESTION_SPACETRACK_CONFIG_FILE="${PWD}/spacetrack.yml"

test $(docker container run -it --rm \
  --volume ${INGESTION_SPACETRACK_VOLUME}:${INGESTION_SPACETRACK_FILE%/*} \
  --workdir ${INGESTION_SPACETRACK_FILE%/*} alpine:latest /bin/sh -c "ls -l" | grep -c "${INGESTION_SPACETRACK_FILE##*/}") -eq 0 && {
    if [[ -z "${INGESTION_SPACETRACK_CONFIG_FILE}" ]]; then
      echo "You need to create the volume ${INGESTION_SPACETRACK_VOLUME} with the following file ${INGESTION_SPACETRACK_FILE}"
      exit 1
    else
      echo "Trying to configure the volume using ${INGESTION_SPACETRACK_CONFIG_FILE} as config file for spacetrack"
      ID=$(docker container run --rm -itd --volume ${INGESTION_SPACETRACK_VOLUME}:${INGESTION_SPACETRACK_FILE%/*} alpine:latest)
      docker container cp ${INGESTION_SPACETRACK_CONFIG_FILE} ${ID}:/home/basket1
      docker container run -it --rm -e "USER=basket1" -e "UID=1001" --volume ingestion_ingsftpvolume:/tmp/upload \
        alpine:latest /bin/sh -c "adduser -h /home/\${USER} -D -u \${UID} -s /sbin/nologin -g \"\" \${USER} && mkdir --parent /tmp/upload/\${USER}/products/automatic && chown -R \${USER}:\${USER} /tmp/upload/\${USER}"

      # Cleanup
      docker container rm -f ${ID}
    fi 
  }

docker container run --name spacetrack-script --detach \
  --volume ingestion_ingsftpvolume:/tmp/upload \
  --volume ${INGESTION_SPACETRACK_VOLUME}:${INGESTION_SPACETRACK_FILE%/*} \
  estenoesmiputonombre/spacetrack:${INGESTION_SPACETRACK_TAG}

sleep 5s

## List all the files inside the target folder where we are downloading spacetrack stuff
ID=$(docker container run --rm -itd --volume ingestion_ingsftpvolume:/tmp/upload alpine:latest /bin/sh)

docker container exec -it ${ID} /bin/sh -c "ls -al /tmp/upload/basket1/products/automatic"
echo "Execute the following command if you want to attach to the script process: \n\"docker container exec --workdir /tmp/upload/basket1/products/automatic -it ${ID}\""

## Copy the file to the current fs
# docker container cp ccc:/tmp/upload/basket1/products/automatic/herethefile ./herethefile

