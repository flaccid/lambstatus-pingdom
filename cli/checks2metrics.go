package main

import (
  "fmt"
  "os"
  "time"

  factory "github.com/flaccid/lambstatus-pingdom/factory"
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
    // option validation
    if len(c.String("lambstatus-endpoint")) < 1 {
      fmt.Println("--lambstatus-endpoint is missing, please provide a value")
      os.Exit(1)
    }
    if len(c.String("lambstatus-api-key")) < 1 {
      fmt.Println("--lambstatus-api-key is missing, please provide a value")
      os.Exit(1)
    }
    if len(c.String("pingdom-user")) < 1 {
      fmt.Println("--pingdom-user is missing, please provide a value")
      os.Exit(1)
    }
    if len(c.String("pingdom-pass")) < 1 {
      fmt.Println("--pingdom-pass is missing, please provide a value")
      os.Exit(1)
    }
    if len(c.String("pingdom-api-key")) < 1 {
      fmt.Println("--pingdom-api-key is missing, please provide a value")
      os.Exit(1)
    }
    if len(c.String("check-to-metric-map")) < 1 {
      fmt.Println("--pingdom-user is missing, please provide a value")
      os.Exit(1)
    }

    factory.Ship(c.String("lambstatus-endpoint"),
         c.String("lambstatus-api-key"),
         c.String("pingdom-user"),
         c.String("pingdom-pass"),
         c.String("pingdom-api-key"),
         c.String("check-to-metric-map"),
         c.Bool("debug"))

    return nil
  }

  app.Run(os.Args)
}
