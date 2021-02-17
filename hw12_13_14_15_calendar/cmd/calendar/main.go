package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/calendar_config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	var rerr error
	defer func() {
		if rerr != nil {
			log.Fatalln(rerr)
		}
	}()

	cfg, err := config.New(configFile)
	if err != nil {
		rerr = err
		return
	}

	logCleanup, err := logger.InitGlobalLogger(cfg)
	if err != nil {
		rerr = err
		return
	}
	defer logCleanup()

	server, cleanup, err := setup(cfg)
	if err != nil {
		rerr = err
		return
	}
	defer cleanup()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	go func() {
		if err := server.GRPC.Start(); err != nil {
			logrus.Warnf("grpc server start failed: %s", err)
			log.Fatalln(err)
		}
	}()

	go func() {
		if err := server.HTTP.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Warnf("http server start failed: %s", err)
			log.Fatalln(err)
		}
	}()

	<-signals
	signal.Stop(signals)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.HTTP.Stop(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Warnf("http server stop failed: %s", err)
	}

	server.GRPC.Stop()
}
