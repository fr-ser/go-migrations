package main

import (
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/fr-ser/go-migrations/commands/start"
	"github.com/fr-ser/go-migrations/utils"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	initLogger()

	app := cli.NewApp()
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		start.StartCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{})

	switch logLevel := utils.GetEnvDefault("LOG_LEVEL", "INFO"); logLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
		log.SetReportCaller(true)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}
