package repositories

import (
	"context"

	"github.com/devesh2997/consequent/identity/domain/entities"
)

type IdentityRepo interface {
	SaveUserPassword(ctx context.Context, userPassword entities.UserPassword) error
	GetActiveUserPassword(ctx context.Context, userID int64) (*entities.UserPassword, error)
	SaveUserLoginMobileOTP(ctx context.Context, otp entities.UserLoginMobileOTP) error
	GetUserLoginMobileOTP(ctx context.Context, verificationID string) (*entities.UserLoginMobileOTP, error)
}
