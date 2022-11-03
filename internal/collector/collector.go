package collector

import (
	"github.com/vukit/magent/internal/config"
	"github.com/vukit/magent/internal/logger"
)

type Collector struct {
	Config *config.Common
	Logger *logger.Logger
}
