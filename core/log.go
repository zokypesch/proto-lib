package core

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"os"

	"github.com/sirupsen/logrus"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmlogrus"
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
	ReportCaller: true,
}

// InitLogWithApm For initialLog
func InitLogWithApm() {
	apm.DefaultTracer.SetLogger(Logs)
	Logs.AddHook(&apmlogrus.Hook{})
}

type Logger struct {
	Logs  *logrus.Logger
	Entry *logrus.Entry
}

func NewLogger() *Logger {
	l := &Logger{Logs: Logs}
	InitLogWithApm()
	return l
}

func (l *Logger) NewLogEntry() *Logger {
	l.Logs.WithFields(logrus.Fields{})

	return l
}

func (l *Logger) LogEntry() *logrus.Entry {
	return l.Entry
}

func (l *Logger) LogWithCtx(ctx context.Context) *Logger {
	traceContextFields := apmlogrus.TraceContext(ctx)
	l.Entry = l.Logs.WithFields(traceContextFields)
	incomingContext, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for k, v := range incomingContext {
			l.Entry = l.Entry.WithField(k, v)
		}
	}

	return l
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

func (l *Logger) LogWithRequest(req interface{}) *Logger {
	l.Entry = l.Entry.WithField("requests_payload", fmt.Sprintf("%s", req))

	return l
}
