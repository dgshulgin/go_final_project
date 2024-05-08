package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func New() *logrus.Logger {
	return &logrus.Logger{
		Out:       os.Stderr,
		Formatter: new(logrus.TextFormatter),
		Level:     logrus.DebugLevel,
	}
}
