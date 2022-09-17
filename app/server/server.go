package server

import (
	"context"
	"net/http"
	"time"

	"github.com/devesh2997/consequent/logger"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, router http.Handler) *Server {
	return &Server{
		httpServer: createServerInstance(port, router),
	}
}

// Serve is...
func (server *Server) Serve(port string, router http.Handler) {
	logger.Log.Debug(context.TODO(), "Server is ready to handle requests at", server.httpServer.Addr)
	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Errorf(context.TODO(), "Could not listen on %s: %v\n", server.httpServer.Addr, err)
	}

}

func createServerInstance(port string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:    port,
		Handler: router,
		// ErrorLog:     logger.Log,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func (server *Server) GracefullyShutdownServer() {
	logger.Log.Debug(context.TODO(), "Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.httpServer.SetKeepAlivesEnabled(false)
	if err := server.httpServer.Shutdown(ctx); err != nil {
		logger.Log.Errorf(context.TODO(), "Could not gracefully shutdown the server: %v\n", err)
	}
}
