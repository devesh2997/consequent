package controllers

import (
	"github.com/devesh2997/consequent/app/controller"
	"github.com/devesh2997/consequent/contextx"
	"github.com/devesh2997/consequent/user/data/mappers"
	"github.com/devesh2997/consequent/user/domain/services"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUser(gCtx *gin.Context)
}

func NewUserController(service services.UserService) UserController {
	return userController{
		service: service,
	}
}

type userController struct {
	controller.Controller
	service services.UserService
}

func (c userController) GetUser(gCtx *gin.Context) {
	requestUser := contextx.GetRequestUser(gCtx.Request.Context())

	user, err := c.service.FindByID(gCtx.Request.Context(), requestUser.ID)
	if err != nil {
		c.SendWithError(gCtx, err)
		return
	}

	userModel := mappers.NewUserMapper().ToModel(*user)

	c.Send(gCtx, userModel)
}
