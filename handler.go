package lambstatus

import (
	"log"
	"net/http"
	"os"

	"github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
	factory "github.com/flaccid/lambstatus-pingdom/factory"
)

var (
	debugMode bool = false
)

func handleEvent() {
	log.Println("function starting")

	factory.Ship(os.Getenv("LAMBSTATUS_ENDPOINT_URL"),
		os.Getenv("LAMBSTATUS_API_KEY"),
		os.Getenv("PINGDOM_USERNAME"),
		os.Getenv("PINGDOM_PASSWORD"),
		os.Getenv("PINGDOM_API_KEY"),
		os.Getenv("CHECK_TO_METRIC_MAP"),
		debugMode)

	log.Println("function complete")
}

// Handle - external handling via runtime context
func Handle(evt interface{}, ctx *runtime.Context) (interface{}, error) {
	handleEvent()
	return "done", nil
}

// HTTPHandle - Cloud Function handler
func HTTPHandle(w http.ResponseWriter, r *http.Request) {
	handleEvent()
}
