package app

import (
	"context"

	"github.com/devesh2997/consequent/config"
	"github.com/devesh2997/consequent/contextx"
	"github.com/devesh2997/consequent/logger"
)

const (
	reqID = "REQUEST_ID"
)

func fieldsExtractor(ctx context.Context) []logger.Field {
	xRequestID := contextx.GetRequestID(ctx)

	f := []logger.Field{
		{
			Key:   reqID,
			Value: xRequestID,
		},
	}

	return f
}

func InitApp(env string) {
	// done := make(chan bool)
	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// config.LoadConfig(env, ".")
	// setupLogger(config.Config.Log)

	// logger := logger.Log
}

func setupLogger(logConfig config.LogConfig) {
	logger.WithFieldsExtractor(logger.LogConfig{
		Level:              logConfig.Level,
		EnableCaller:       logConfig.EnableCaller,
		ErrOutputFilePaths: []string{"./storage/logs/app.log"},
		OutputFilePaths:    []string{"./storage/logs/app.log"},
	}, fieldsExtractor)
}
