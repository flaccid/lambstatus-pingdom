package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/russellcardullo/go-pingdom/pingdom"
)

func test() {
	fmt.Println("LAMBSTATUS_ENDPOINT_URL=" + os.Getenv("LAMBSTATUS_ENDPOINT_URL"))
	fmt.Println("CHECK_TO_METRIC_MAP=" + os.Getenv("CHECK_TO_METRIC_MAP"))

	pClient := pingdom.NewClient(os.Getenv("PINGDOM_USERNAME"),
		os.Getenv("PINGDOM_PASSWORD"),
		os.Getenv("PINGDOM_API_KEY"))

	checkToMetricMap := os.Getenv("CHECK_TO_METRIC_MAP")

	// split the map by comma
	mappings := strings.Split(checkToMetricMap, ",")
	fmt.Println("Processing", len(mappings), "mappings")

	// iterate through each mapping
	for i := range mappings {
		mapping := strings.Split(mappings[i], ":")

		fmt.Println("\r\nPingdom Check ID:", mapping[0])
		fmt.Println("LambStatus Metric ID:", mapping[1])

		checkId, err := strconv.Atoi(mapping[0])
		metricId := mapping[1]

		// get details for the check
		checkDetails, _ := pClient.Checks.Read(checkId)
		fmt.Printf("Check Details: %+v\n", checkDetails)
		lastResponseTime := checkDetails.LastResponseTime
		lastTestTime := checkDetails.LastTestTime

		fmt.Println(fmt.Sprintf("Last Response Time: %vms", lastResponseTime))
		fmt.Println(fmt.Sprintf("Last Test Time: %v ", lastTestTime))

		t := time.Unix(lastTestTime, 0).Add(1).UTC().Format(time.RFC3339Nano)

		var jsonPayLoad = []byte("{\"" + metricId + `": [{"timestamp": "` + t + `", "value": ` + fmt.Sprintf("%v", lastResponseTime) + `}]}`)
		fmt.Println("JSON payload:", string(jsonPayLoad[:]))

		url := os.Getenv("LAMBSTATUS_ENDPOINT_URL") + "/prod/v0/metrics/data"
		fmt.Println("POST:", url)

		// send to lambstatus
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayLoad))
		req.Header.Set("x-api-key", os.Getenv("LAMBSTATUS_API_KEY"))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}
}
