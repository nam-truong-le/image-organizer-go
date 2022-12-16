package log

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	logger     *logrus.Logger
	initLogger sync.Once
)

func Logger() *logrus.Logger {
	initLogger.Do(func() {
		logger = logrus.New()
	})
	return logger
}
