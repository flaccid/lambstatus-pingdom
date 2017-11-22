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
