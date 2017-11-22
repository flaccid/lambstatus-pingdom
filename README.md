[![CircleCI](https://circleci.com/gh/flaccid/lambstatus-pingdom.svg?style=svg)](https://circleci.com/gh/flaccid/lambstatus-pingdom)

# lambstatus-pingdom

Ship pingdom to lambstatus.

## Usage

### Go

This basic example is kept in `contrib/example.go`:

```
package main

import (
	"os"

	factory "github.com/flaccid/lambstatus-pingdom/factory"
)

var (
	debugMode bool = false
)

func main() {
	factory.Ship(os.Getenv("LAMBSTATUS_ENDPOINT_URL"),
		os.Getenv("LAMBSTATUS_API_KEY"),
		os.Getenv("PINGDOM_USERNAME"),
		os.Getenv("PINGDOM_PASSWORD"),
		os.Getenv("PINGDOM_API_KEY"),
		os.Getenv("CHECK_TO_METRIC_MAP"),
		debugMode)
}
```

### CLI

Using the `checks2metrics` command, you can provide CLI options or just set the environment variables below.

#### Environment Variables

- `PINGDOM_USERNAME` - pingdom username
- `PINGDOM_PASSWORD` - pingdom password
- `PINGDOM_API_KEY` - pingdom API key
- `LAMBSTATUS_ENDPOINT_URL` - LambStatus endpoint URL
- `LAMBSTATUS_API_KEY` - LambStatus API key
- `CHECK_TO_METRIC_MAP` - map of pingdom checks to lambstatus metrics;
in format, `<pingdom_check_id>:<lambstatus_metric_id>,<pingdom_check_id>:<lambstatus_metric_id>`

Store your environment variables in `.env`.

    $ source .env && go run cli/checks2metrics.go

### AWS Lambda

#### Setup

First, create the `lambda_basic_execution` IAM role for use with the function:

    $ make lambda-create-iam-role

Then, attach the `arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole` policy:

    $ make lambda-attach-iam-role-policy

Now, create the `checks2metrics` function itself:

    $ make lambda-create-function

Build the `handler.zip`:

    $ make lambda-build-pack

Upload the pack:

    $ make lambda-update-function-code

#### Running the Function

Note: you may need to first increase the function timeout.

    $ make lambda-invoke-function

#### Scheduling the Function

In reality, its most likely you'd like to run this every 1-2 minutes using http://docs.aws.amazon.com/lambda/latest/dg/with-scheduled-events.html.

You can automatically create a 1 minute schedule using (this is still a bit experimental):

		$ make aws-create-scheduled-event

### Building

#### CLI

Standard local build:

    $ go build -o bin/checks2metrics checks2metrics.go

Static 64-bit Linux executable:

    $ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
        go build \
        -o bin/checks2metrics \
        -a -ldflags '-extldflags "-static"' \
          checks2metrics.go

#### Docker Image

    $ docker build -t flaccid/checks2metrics .

Push to Docker Hub:

    $ docker push flaccid/checks2metrics

License and Authors
-------------------
- Author: Chris Fordham (<chris@fordham-nagy.id.au>)

```text
Copyright 2017, Chris Fordham

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
