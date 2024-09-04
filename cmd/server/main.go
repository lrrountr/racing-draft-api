package main

import (
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/lrrountr/racing-draft-api/internal/config"
	"github.com/lrrountr/racing-draft-api/internal/handler"
)

const (
	APP_RETRY_TIMER = time.Second
)

func main() {
	log.SetReportCaller(true)
	app := getApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func getApp() *cli.App {
	return &cli.App{
		Name:     "racing-draft",
		Usage:    "Start Racing App server",
		Action:   StartServer,
		Flags:    []cli.Flag{},
		Commands: []*cli.Command{},
	}
}

func StartServer(c *cli.Context) error {
	config, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("Failed to load config from environment: %w", err)
	}
	for {
		err := startServer(config)
		if err != nil {
			log.Printf("%s\n", err)
			time.Sleep(APP_RETRY_TIMER)
		}
	}
}

func startServer(config config.Config) error {
	return handler.StartServer(config)
}
