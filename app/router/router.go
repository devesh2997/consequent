package router

import (
	"net/http"

	"github.com/devesh2997/consequent/app/middleware"
	"github.com/devesh2997/consequent/logger"
	"github.com/gin-gonic/gin"
)

// Create is...
func Create() http.Handler {
	r := gin.New()

	setupGlobalMiddlewares(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello from consequent.",
		})
	})

	return r
}

func setupGlobalMiddlewares(r *gin.Engine) {
	accessLogWriter, err := logger.NewAccessLogWriter("storage/logs/access.log")
	if err != nil {
		panic(err)
	}

	r.Use(gin.LoggerWithWriter(accessLogWriter))

	r.Use(gin.RecoveryWithWriter(logger.NewErrorLogWriter()))
	r.Use(middleware.RequestInfo(nil))
}
