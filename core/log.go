package core

import (
	"os"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
)

// Logs for setup logs
var Logs = &logrus.Logger{
	Out:   os.Stderr,
	Hooks: make(logrus.LevelHooks),
	Level: logrus.DebugLevel,
	Formatter: &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "log.level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "function.name", // non-ECS
		},
	},
}

// InitLogWithApm For initialLog
func InitLogWithApm() {
	apm.DefaultTracer.SetLogger(Logs)
}
