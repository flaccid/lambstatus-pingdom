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
