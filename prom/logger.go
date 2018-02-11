package prom

import (
	"os"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var logger = &logrus.Logger{
	Out:       os.Stdout,
	Formatter: &prefixed.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", ForceColors: true, ForceFormatting: true, FullTimestamp: true},
	Hooks:     make(logrus.LevelHooks),
	Level:     logrus.InfoLevel,
}
