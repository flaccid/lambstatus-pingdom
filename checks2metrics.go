package main

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "net/http"
  log "github.com/sirupsen/logrus"
  "os"
  "strconv"
  "strings"
  "time"

  "github.com/russellcardullo/go-pingdom/pingdom"
  "github.com/urfave/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "checks2metrics"
  app.Version = "v0.0.1"
  app.Compiled = time.Now()
  app.Copyright = "(c) 2016 Chris Fordham"
  app.Authors = []cli.Author{
		cli.Author{
			Name:  "Chris Fordham",
			Email: "chris@fordham-nagy.id.au",
		},
	}
  app.Usage = "Ship Pingdom checks to LambStatus metrics."
  app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pingdom-user",
      EnvVar: "PINGDOM_USERNAME",
		},
		cli.StringFlag{
			Name:  "pingdom-pass",
      EnvVar: "PINGDOM_PASSWORD",
		},
    cli.StringFlag{
			Name:  "pingdom-api-key",
      EnvVar: "PINGDOM_API_KEY",
		},
    cli.StringFlag{
			Name:  "lambstatus-endpoint",
      EnvVar: "LAMBSTATUS_ENDPOINT_URL",
		},
    cli.StringFlag{
			Name:  "lambstatus-api-key",
      EnvVar: "LAMBSTATUS_API_KEY",
		},
    cli.StringFlag{
      Name:  "check-to-metric-map",
      EnvVar: "CHECK_TO_METRIC_MAP",
    },
    cli.BoolFlag{
      Name:  "debug",
      Usage: "set debug log level",
    },
	}

  app.Action = func(c *cli.Context) error {
    log.WithFields(log.Fields{
      "lambstatus endpoint url": c.String("lambstatus-endpoint"),
      "check2metric map": c.String("check-to-metric-map"),
    }).Info("starting")

    if c.Bool("debug") {
      log.Info("debug enabled")
      log.SetLevel(log.DebugLevel)
    }

    pClient := pingdom.NewClient(c.String("pingdom-user"),
                                 c.String("pingdom-pass"),
                                 c.String("pingdom-api-key"))
    checkToMetricMap := c.String("check-to-metric-map")

    // split the map by comma
    mappings := strings.Split(checkToMetricMap, ",")
    log.Info(len(mappings), " checks to ship")

    // iterate through each mapping
    for i := range mappings {
      mapping := strings.Split(mappings[i], ":")
      checkId, err := strconv.Atoi(mapping[0])
      metricId := mapping[1]
      checkDetails, err := pClient.Checks.Read(checkId)
      if err != nil {
        log.Fatal(err)
      }
      log.Debug("check details: %+v\n", checkDetails)
      lastResponseTime := checkDetails.LastResponseTime
      lastTestTime := checkDetails.LastTestTime

      dateStamp := time.Unix(lastTestTime, 0).Add(1).UTC().Format(time.RFC3339Nano)
      var jsonPayLoad = []byte("{\""+metricId+`": [{"timestamp": "` + dateStamp + `", "value": ` + fmt.Sprintf("%v", lastResponseTime) + `}]}`)

      log.WithFields(log.Fields{
        "pingdom check id": checkId,
        "lambstatus metric id": metricId,
        "last response time": lastResponseTime,
        "last test time": lastTestTime,
        "timestamp": dateStamp,
      }).Info("ship ", i)
      log.Debug("JSON payload:", string(jsonPayLoad[:]))

      url := c.String("lambstatus-endpoint") + "/prod/v0/metrics/data"
      log.Debug("POST: ", url)

      // send to lambstatus
      req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayLoad))
      req.Header.Set("x-api-key", c.String("lambstatus-api-key"))
      req.Header.Set("Content-Type", "application/json")

      client := &http.Client{}
      resp, err := client.Do(req)
      if err != nil {
        panic(err)
      }
      defer resp.Body.Close()

      log.Debug("response: ", resp.Status)
      log.Debug("response headers: ", resp.Header)
      body, _ := ioutil.ReadAll(resp.Body)
      log.Debug("response body: ", string(body))
    }
    log.Info("all checks shipped")

    return nil
  }

  app.Run(os.Args)
}
