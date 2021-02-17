package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/scheduler/scheduler_config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	scheduler, cleanup, err := setup(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer cleanup()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	go func() {
		if err := scheduler.Run(context.Background()); err != nil {
			logrus.WithError(err).Error("scheduler run failed")
			log.Fatalln(err)
		}
	}()

	<-signals
	signal.Stop(signals)

	if err := scheduler.Shutdown(); err != nil {
		logrus.WithError(err).Error("shutdown failed")
	}
}
