package services

import (
	"context"

	"github.com/devesh2997/consequent/identity/domain/entities"
	"github.com/devesh2997/consequent/user/domain/services"
)

type IdentityService interface {
	SendOTP(ctx context.Context, mobileNumber int64) error
	VerifyOTP(ctx context.Context, mobileNumber int64, otp int) (*entities.Token, error)
}

func NewIdentityService(userService services.UserService) IdentityService {
	return identityService{userService: userService}
}

type identityService struct {
	userService services.UserService
	otpSender   OTPSender
}

func (identityService) SendOTP(ctx context.Context, mobileNumber int64) error {
	panic("unimplemented")
}

func (identityService) VerifyOTP(ctx context.Context, mobileNumber int64, otp int) (*entities.Token, error) {
	panic("unimplemented")
}
