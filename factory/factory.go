package factory

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/russellcardullo/go-pingdom/pingdom"
)

func Ship(lambStatusEndpoint string,
	lambStatusApiKey string,
	pingdomUser string,
	pingdomPass string,
	pingdomApiKey string,
	checkToMetricMap string,
	debugMode bool) {

	log.WithFields(log.Fields{
		"lambstatus endpoint url": lambStatusEndpoint,
		"checks to metrics map":   checkToMetricMap,
	}).Info("starting")

	if debugMode {
		log.Info("debug enabled")
		log.SetLevel(log.DebugLevel)
	}

	pClient := pingdom.NewClient(pingdomUser, pingdomPass, pingdomApiKey)

	// split the map by comma
	mappings := strings.Split(checkToMetricMap, ",")
	log.Info(len(mappings), " checks to ship to lambstatus")

	// iterate through each mapping
	for i := range mappings {
		mapping := strings.Split(mappings[i], ":")
		checkId, err := strconv.Atoi(mapping[0])
		metricId := mapping[1]
		checkDetails, err := pClient.Checks.Read(checkId)
		if err != nil {
			// expect non-200 responses or complete fails
			log.WithFields(log.Fields{
				"pingdom check id":     checkId,
				"lambstatus metric id": metricId,
				"error":								err,
				"index":								i,
			}).Error("failure getting pingdom check")
		} else {
			log.Debug("check details: %+v\n", checkDetails)
			lastResponseTime := checkDetails.LastResponseTime
			lastTestTime := checkDetails.LastTestTime

			dateStamp := time.Unix(lastTestTime, 0).Add(1).UTC().Format(time.RFC3339Nano)
			var jsonPayLoad = []byte("{\"" + metricId + `": [{"timestamp": "` + dateStamp + `", "value": ` + fmt.Sprintf("%v", lastResponseTime) + `}]}`)

			log.WithFields(log.Fields{
				"pingdom check id":     checkId,
				"lambstatus metric id": metricId,
				"last response time":   lastResponseTime,
				"last test time":       lastTestTime,
				"datestamp":            dateStamp,
			}).Info(i, " ", checkDetails.Name)
			log.Debug("JSON payload:", string(jsonPayLoad[:]))

			url := lambStatusEndpoint + "/api/v0/metrics/data"
			log.Debug("POST: ", url)

			// send to lambstatus
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayLoad))
			req.Header.Set("x-api-key", lambStatusApiKey)
			req.Header.Set("Content-Type", "application/json")
			log.Debug("request", req)

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				// expect non-200 responses or complete fails
				log.Errorf("failure sending metric to lambstatus: %s", err)
			} else {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err)
				}	else {
					log.Debug("response status: ", resp.Status)
					log.Debug("response headers: ", resp.Header)
					log.Debug("response body: ", string(body))
				}
				if resp.StatusCode != 200 {
					log.Error("failed to post metric: ", resp.Status, " ", string(body))
				}
			}
		}
	}
	log.Info("all shippable checks shipped")
}
