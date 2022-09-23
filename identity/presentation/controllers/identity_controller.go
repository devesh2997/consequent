package controllers

import (
	"errors"

	"github.com/devesh2997/consequent/app/controller"
	"github.com/devesh2997/consequent/identity/data/mappers"
	"github.com/devesh2997/consequent/identity/domain/services"
	"github.com/gin-gonic/gin"
)

type IdentityController interface {
	SendOTP(gCtx *gin.Context)
	VerifyOTP(gCtx *gin.Context)
	IsEmailRegistered(gCtx *gin.Context)
	SignUpWithEmail(gCtx *gin.Context)
	SignInWithEmailAndPassword(gCtx *gin.Context)
}

func NewIdentityController(service services.IdentityService) IdentityController {
	return identityController{service: service}
}

type identityController struct {
	controller.Controller
	service services.IdentityService
}

func (c identityController) SendOTP(gCtx *gin.Context) {
	mobileNumber, exists := gCtx.GetQuery("mobile_number")
	if !exists {
		c.SendBadRequestError(gCtx, errors.New("mobile number is required"))
		return
	}

	verificationID, err := c.service.SendOTP(gCtx.Request.Context(), mobileNumber)
	if err != nil {
		c.SendWithError(gCtx, err)
		return
	}

	c.Send(gCtx, gin.H{
		"verification_id": verificationID,
	})
}

func (c identityController) VerifyOTP(gCtx *gin.Context) {
	input := struct {
		VerificationID string `json:"verification_id" form:"verification_id"`
		MobileNumber   string `json:"mobile_number" form:"mobile_number"`
		OTP            int    `json:"otp" form:"otp"`
	}{}

	if err := gCtx.ShouldBind(&input); err != nil {
		c.SendBadRequestError(gCtx, err)
		return
	}

	token, err := c.service.VerifyOTP(gCtx.Request.Context(), input.VerificationID, input.MobileNumber, input.OTP)
	if err != nil {
		c.SendWithError(gCtx, err)
		return
	}

	tokenModel := mappers.NewTokenMapper().ToModel(*token)

	c.Send(gCtx, tokenModel)
}

func (c identityController) IsEmailRegistered(gCtx *gin.Context) {
	email, exists := gCtx.GetQuery("email")
	if !exists {
		c.SendBadRequestError(gCtx, errors.New("email is required"))
		return
	}
	registered, err := c.service.IsEmailRegistered(gCtx.Request.Context(), email)
	if err != nil {
		c.SendWithError(gCtx, err)
		return
	}

	c.Send(gCtx, gin.H{
		"registered": registered,
	})
}

func (c identityController) SignUpWithEmail(gCtx *gin.Context) {
	input := struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}{}

	if err := gCtx.ShouldBind(&input); err != nil {
		c.SendBadRequestError(gCtx, err)
		return
	}

	token, err := c.service.SignUpWithEmail(gCtx.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.SendWithError(gCtx, err)
		return
	}

	tokenModel := mappers.NewTokenMapper().ToModel(*token)

	c.Send(gCtx, tokenModel)
}

func (c identityController) SignInWithEmailAndPassword(gCtx *gin.Context) {
	input := struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}{}

	if err := gCtx.ShouldBind(&input); err != nil {
		c.SendBadRequestError(gCtx, err)
		return
	}

	token, err := c.service.SignInWithEmailAndPassword(gCtx.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.SendWithError(gCtx, err)
		return
	}

	tokenModel := mappers.NewTokenMapper().ToModel(*token)

	c.Send(gCtx, tokenModel)
}
