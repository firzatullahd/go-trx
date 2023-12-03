package logger

import (
	"context"
	log "github.com/sirupsen/logrus"
	"go-trx/config"
	"os"
)

func Init(conf config.Config) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
}

func Info(ctx context.Context, format string, values ...interface{}) {
	var id string
	val := ctx.Value("X-Correlation-ID")
	if val != nil {
		id = val.(string)
	}
	log.WithFields(log.Fields{
		"Correlation-ID": id,
	}).Infof(format, values...)
}

func Warn(ctx context.Context, format string, values ...interface{}) {
	var id string
	val := ctx.Value("X-Correlation-ID")
	if val != nil {
		id = val.(string)
	}
	log.WithFields(log.Fields{
		"Correlation-ID": id,
	}).Warnf(format, values...)
}

func Error(ctx context.Context, format string, values ...interface{}) {
	var id string
	val := ctx.Value("X-Correlation-ID")
	if val != nil {
		id = val.(string)
	}
	log.WithFields(log.Fields{
		"Correlation-ID": id,
	}).Errorf(format, values...)
}

func Debug(ctx context.Context, format string, values ...interface{}) {
	var id string
	val := ctx.Value("X-Correlation-ID")
	if val != nil {
		id = val.(string)
	}
	log.WithFields(log.Fields{
		"Correlation-ID": id,
	}).Debugf(format, values...)
}

func Fatal(ctx context.Context, format string, values ...interface{}) {
	var id string
	val := ctx.Value("X-Correlation-ID")
	if val != nil {
		id = val.(string)
	}
	log.WithFields(log.Fields{
		"Correlation-ID": id,
	}).Fatalf(format, values...)
}
