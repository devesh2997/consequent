package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/devesh2997/consequent/app/router"
	"github.com/devesh2997/consequent/app/server"
	"github.com/devesh2997/consequent/config"
	"github.com/devesh2997/consequent/contextx"
	"github.com/devesh2997/consequent/datasources"
	"github.com/devesh2997/consequent/logger"
)

const (
	reqID       = "REQUEST_ID"
	defaultPort = ":5035"
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
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	config.LoadConfig(env, ".")
	setupLogger(config.Config.Log)

	logger := logger.Log

	ds, err := datasources.Get()

	if err != nil {
		logger.Debug(context.TODO(), err.Error())
		return
	}

	// Setup http server
	router := router.Create()

	logger.Debug(context.TODO(), "Server is starting...")

	port := config.Config.Port
	if port == "" {
		port = defaultPort
		fmt.Printf("Connected to default port: %s", port)
	}

	httpServer := server.NewServer(port, router)
	go func() {
		<-quit
		httpServer.GracefullyShutdownServer()
		close(done)
	}()

	httpServer.Serve(port, router)

	<-done

	ds.SQLClients.GetSQLXClusterDB().Close()
	logger.Debug(context.TODO(), "Server stopped")
}

func setupLogger(logConfig config.LogConfig) {
	logger.WithFieldsExtractor(logger.LogConfig{
		Level:              logConfig.Level,
		EnableCaller:       logConfig.EnableCaller,
		ErrOutputFilePaths: []string{"./storage/logs/app.log"},
		OutputFilePaths:    []string{"./storage/logs/app.log"},
	}, fieldsExtractor)
}
