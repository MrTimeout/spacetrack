# Space track script

This script is used to fetch data about man-made earth-orbiting object tracked by [www.space-track.org](https://www.space-track.org).

## Execute

### Docker

#### No Secure

This mode is used when we don't need/want to use a passphrase to encrypt our credentials.

```sh
> docker image build -f ./Dockerfile.nosecure -t spacetrack:latest -t spacetrack:0.0.1-SNAPSHOT --no-cache --target nosecure .
> docker container run --name spacetrack \
	--rm \
	--detach \
	-u spacetrack-client:spacetrack-client \
	--mount "type=volume,source=spacetrack,target=/tmp" \
	--mount "type=bind,source=$HOME/.spacetrack.yaml,target=/home/spacetrack-client/.spacetrack.yaml" \
	spacetrack:latest
```

#### Secure

This mode is used when we want to encrypt our credentials. We need to use an orchestrator, this example is from docker swarm. The orchestrator used, is up to you.

```sh
> docker image build -f ./Dockerfile -t spacetrack:latest -t spacetrack:0.0.1-SNAPSHOT --no-cache --target pre .
> docker secret create spacetrack_secret ./secret-file
> docker config create spacetrack_config ~/.spacetrack.yaml
> docker stack deploy --compose-file ./docker-compose.yml spacetrack
```

## Commands

### GP (General Perturbations)

The general perturbations (GP) class is an efficient listing of the newest SGP4 keplerian element set for each man-made earth-orbiting object tracked by the 18th Space Defense Squadron.

#### Flags

- `--format`: format flag allows us to get the data in different formats:
	+ xml
	+ csv
	+ html
	+ json
- `--sort`: Sort response Ascending or Descending. By default, it is asc:
	+ asc
	+ desc
- `--limit`: Limitting output to a restrictive number of results
- `--skip`: Skipping first n elements.
- `--dry-run`: Execute script and return the target path only without executing REST.
- `--orderby`: Order results by specified field, which is present on the response. Posible fields(lower or upper case):
	+ CCSDS_OMM_VERS
	+ COMMENT
	+ CREATION_DATE
	+ ORIGINATOR
	+ OBJECT_NAME
	+ OBJECT_ID
	+ CENTER_NAME
	+ REF_FRAME
	+ TIME_SYSTEM
	+ MEAN_ELEMENT_THEORY
	+ EPOCH
	+ MEAN_MOTION
	+ ECCENTRICITY
	+ INCLINATION
	+ RA_OF_ASC_NODE
	+ ARG_OF_PERICENTER
	+ MEAN_ANOMALY
	+ EPHEMERIS_TYPE
	+ CLASSIFICATION_TYPE
	+ NORAD_CAT_ID
	+ ELEMENT_SET_NO
	+ REV_AT_EPOCH
	+ BSTAR
	+ MEAN_MOTION_DOT
	+ MEAN_MOTION_DDOT
	+ SEMIMAJOR_AXIS
	+ PERIOD
	+ APOAPSIS
	+ PERIAPSIS
	+ OBJECT_TYPE
	+ RCS_SIZE
	+ COUNTRY_CODE
	+ LAUNCH_DATE
	+ SITE
	+ DECAY_DATE
	+ FILE
	+ GP_ID
	+ TLE_LINE0
	+ TLE_LINE1
	+ TLE_LINE2
- `--filter`: We can filter by the same fields as orderby. We can use different operators:
	+ `>`: More than
	+ `<`: Less than
	+ `=`: Equal
	+ `<>`: Not equal
	+ `~~`: Contains
	+ `^`: Starts with
	+ `--`: Range of numbers

### Examples

`spacetrack gp --filter 'decay_date<>null-val' --filter 'OBJECT_ID=1960,1961--1964,^196,=1970' --orderby NORAD_CAT_ID --sort asc --limit 5 --skip 4`

- Filter by `DECAY_DATE` not equal null and `OBJECT_ID` be between 1960 and 1970 (which represents the year of launch date).
- Order by `NORAD_CAT_ID` which is a number that represents the orbital object
- Sort in Ascending order
- Limit the response to 5 rows and skip the first 4

## Flags

- `-i/--interval`: Interval in `time.Duration` format. Minimum value is 5 minutes. By default is 30 minutes.
- `-w/--work-dir`: Folder where all files are going to be persisted. It is requried.
- `-lf/--log-file`: File where all the logs are going to. If it is not already created, it will be created for you. If it already exists, the logs will be appended to the file. It is not required. By default, logs are printed to the output console.
- `-l/--log-level`: It allow us to set the logging level of the application, it could be:
	+ debug or DEBUG
	+ info or INFO
	+ warn or WARN
	+ error or ERROR
	+ fatal or FATAL
