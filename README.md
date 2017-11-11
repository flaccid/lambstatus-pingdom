# lambstatus-pingdom

Ship pingdom to lambstatus.

## Usage

Store your environment variables in `.env`.

    $ source .env && go run test.go

### Environment Variables

- `PINGDOM_USERNAME`
- `PINGDOM_PASSWORD`
- `PINGDOM_API_KEY`
- `LAMBSTATUS_ENDPOINT_URL`
- `LAMBSTATUS_API_KEY`
- `CHECK_TO_METRIC_MAP`

In format, `<pingdom_check_id>:<lambstatus_metric_id>,<pingdom_check_id>:<lambstatus_metric_id>`.

## Building

Standard local build:

    $ go build -o bin/checks2metrics checks2metrics.go

Static 64-bit Linux executable:

    $ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build \
        -o bin/checks2metrics \
        -a -ldflags '-extldflags "-static"' \
          checks2metrics.go

### Docker Image

    $ docker build -t flaccid/checks2metrics .

Push to Docker Hub:

    $ docker push flaccid/checks2metrics
