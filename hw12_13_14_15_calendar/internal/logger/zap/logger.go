package zap

import (
	"fmt"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

var level = map[string]zapcore.Level{
	"info":    zapcore.InfoLevel,
	"error":   zapcore.ErrorLevel,
	"warning": zapcore.WarnLevel,
	"debug":   zapcore.DebugLevel,
}

func New(cfg *config.Config) (*Logger, error) {
	if _, ok := level[cfg.Logger.Level]; !ok {
		return nil, fmt.Errorf("unexpected logger level %s", cfg.Logger.Level)
	}

	lcfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(level[cfg.Logger.Level]),
		OutputPaths: []string{cfg.Logger.Path},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "time",
			EncodeTime:  zapcore.RFC3339TimeEncoder,
			NameKey:     "name",
			EncodeName:  zapcore.FullNameEncoder,
		},
	}

	zapLogger, err := lcfg.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger failed: %w", err)
	}

	return &Logger{
		logger: zapLogger,
	}, nil
}

func (zl *Logger) Named(name string) logger.Logger {
	return &Logger{
		logger: zl.logger.Named(name),
	}
}

func (zl *Logger) Infof(msg string, vals ...interface{}) {
	zl.logger.Info(fmt.Sprintf(msg, vals...))
}

func (zl *Logger) Errorf(msg string, vals ...interface{}) {
	zl.logger.Error(fmt.Sprintf(msg, vals...))
}

func (zl *Logger) Warnf(msg string, vals ...interface{}) {
	zl.logger.Warn(fmt.Sprintf(msg, vals...))
}

func (zl *Logger) Debugf(msg string, vals ...interface{}) {
	zl.logger.Debug(fmt.Sprintf(msg, vals...))
}
