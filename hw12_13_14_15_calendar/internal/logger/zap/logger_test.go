package zap

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sterligov/otus_homework/hw12_13_14_15_calendar/internal/config"
)

func TestLogger(t *testing.T) {
	levels := []string{"info", "error", "warning", "debug"}

	var cfg config.Config
	cfg.Logger.Path = "stderr"

	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			cfg.Logger.Level = level
			logger, err := New(&cfg)

			require.NoError(t, err)
			require.NotNil(t, logger)
		})
	}

	t.Run("unexpected", func(t *testing.T) {
		cfg.Logger.Level = "unexpected"
		logger, err := New(&cfg)

		require.Error(t, err)
		require.Nil(t, logger)
	})
}
