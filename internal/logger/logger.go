package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"

	"github.com/vukit/magent/internal/config"
)

type Logger struct {
	logger zerolog.Logger
}

func NewLogger(mConfig *config.Config, w io.Writer) *Logger {
	r := &Logger{}

	if mConfig.Common.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if w == nil {
		w = os.Stderr
	}

	r.logger = zerolog.New(w).With().Timestamp().Str("host", mConfig.Common.HostName).Logger()

	return r
}

func (r *Logger) Debug(message interface{}) {
	r.logger.Debug().Interface("message", message).Send()
}

func (r *Logger) Fatal(message interface{}) {
	r.logger.Fatal().Interface("message", message).Send()
}

func (r *Logger) Warning(message interface{}) {
	r.logger.Warn().Interface("message", message).Send()
}

func (r *Logger) Info(message interface{}) {
	r.logger.Info().Interface("message", message).Send()
}
