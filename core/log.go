package core

import (
	"context"
	"google.golang.org/grpc/metadata"
	"os"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
)

// Logs for setup logs
var Logs = &logrus.Logger{
	Out:   os.Stdout,
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
	ReportCaller: true,
}

// InitLogWithApm For initialLog
func InitLogWithApm() {
	apm.DefaultTracer.SetLogger(Logs)
	Logs.AddHook(&apmlogrus.Hook{})
}

func AddCtxToLog(ctx context.Context, logger *logrus.Logger) *logrus.Entry {
	withFields := logger.WithFields(logrus.Fields{})
	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for k, v := range incomingContext {
			withFields = logger.WithField(k, v)
		}
	}

	return withFields
}
