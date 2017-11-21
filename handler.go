package main

import (
	"log"
	"os"

	factory "github.com/flaccid/lambstatus-pingdom/factory"
	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
)

var (
  debugMode bool = false
)

func Handle(evt interface{}, ctx *runtime.Context) (interface{}, error) {
	log.Println("function starting")

	factory.Ship(os.Getenv("LAMBSTATUS_ENDPOINT_URL"),
			         os.Getenv("LAMBSTATUS_API_KEY"),
			         os.Getenv("PINGDOM_USERNAME"),
			         os.Getenv("PINGDOM_PASSWORD"),
			         os.Getenv("PINGDOM_API_KEY"),
			         os.Getenv("CHECK_TO_METRIC_MAP"),
               debugMode)

	log.Println("function complete")

	return "done", nil
}
