package entities

import "time"

type UserLoginMobileOTP struct {
	ID             int64
	VerificationID string
	Mobile         int64
	OTP            int
	Status         string
	CreatedAt      time.Time
	ExpiryAt       time.Time
	UpdatedAt      time.Time
}

func (otp UserLoginMobileOTP) HasExpired() bool {
	return time.Now().After(otp.ExpiryAt)
}
