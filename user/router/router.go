package router

import (
	"github.com/devesh2997/consequent/app/middleware"
	identityContainers "github.com/devesh2997/consequent/identity/containers"
	"github.com/devesh2997/consequent/user/containers"
	"github.com/gin-gonic/gin"
)

func InjectUserRoutes(router *gin.RouterGroup) {
	setupV1Routes(router)
}

func setupV1Routes(r *gin.RouterGroup) {
	tokenService := identityContainers.InjectTokenService()
	userController := containers.InjectUserController()

	v1 := r.Group("/v1")
	v1.Use(middleware.Authorisation(tokenService))
	v1.GET("user", func(c *gin.Context) {
		userController.GetUser(c)
	})
}
