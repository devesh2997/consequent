package entities

import (
	"time"

	"github.com/devesh2997/consequent/identity/constants"
)

type UserLoginMobileOTP struct {
	ID             int64
	VerificationID string
	Mobile         string
	OTP            int
	Status         string
	CreatedAt      time.Time
	ExpiryAt       time.Time
	UpdatedAt      time.Time
}

func (otp UserLoginMobileOTP) IsActive() bool {
	return otp.Status == constants.USER_LOGIN_MOBILE_OTP_STATUS_ACTIVE
}

func (otp UserLoginMobileOTP) HasExpired() bool {
	return time.Now().After(otp.ExpiryAt)
}
