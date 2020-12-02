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

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/logger"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/logger/zap"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := config.New(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	zlog, err := zap.New(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	logger.SetGlobalLogger(zlog)

	var exitCode int

	server, cleanup, err := setup(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		cleanup()
		os.Exit(exitCode)
	}()

	go func() {
		if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("http server start failed: %s", err)
			exitCode = 1
			return
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)

	<-signals
	signal.Stop(signals)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Stop(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("http server stop failed: %s", err)
		exitCode = 1
		return
	}
}
