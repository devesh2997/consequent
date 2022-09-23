package router

import (
	"github.com/devesh2997/consequent/identity/containers"
	"github.com/gin-gonic/gin"
)

func InjectIdentityRoutes(router *gin.RouterGroup) {
	setupV1Routes(router)
}

func setupV1Routes(r *gin.RouterGroup) {
	identiyController := containers.InjectIdentityController()

	v1 := r.Group("/v1")
	v1.POST("/send-otp", func(c *gin.Context) {
		identiyController.SendOTP(c)
	})
	v1.POST("/verify-otp", func(c *gin.Context) {
		identiyController.VerifyOTP(c)
	})
	v1.GET("/is-email-registered", func(c *gin.Context) {
		identiyController.IsEmailRegistered(c)
	})
	v1.POST("/sign-up-with-email", func(c *gin.Context) {
		identiyController.SignUpWithEmail(c)
	})
	v1.POST("/sign-in-with-email", func(c *gin.Context) {
		identiyController.SignInWithEmailAndPassword(c)
	})
}
