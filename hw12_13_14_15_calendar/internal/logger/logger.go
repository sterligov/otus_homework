package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
)

var level = map[string]logrus.Level{
	"info":    logrus.InfoLevel,
	"error":   logrus.ErrorLevel,
	"warning": logrus.WarnLevel,
	"debug":   logrus.DebugLevel,
}

func InitGlobalLogger(cfg *config.Config) (func(), error) {
	logClose := func() {}

	if _, ok := level[cfg.Logger.Level]; !ok {
		return logClose, fmt.Errorf("unexpected logger level %s", cfg.Logger.Level)
	}

	var (
		flog *os.File
		err  error
	)

	switch cfg.Logger.Path {
	case "stderr":
		flog = os.Stderr
	case "stdout":
		flog = os.Stdout
	default:
		flog, err = os.OpenFile(cfg.Logger.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0755))
		if err != nil {
			return logClose, fmt.Errorf("open log file failed: %w", err)
		}
	}

	logrus.SetOutput(flog)
	logrus.SetLevel(level[cfg.Logger.Level])
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logClose = func() {
		if err := flog.Close(); err != nil {
			logrus.Warn(flog)
		}
	}

	return logClose, nil
}
