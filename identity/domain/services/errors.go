package services

import "github.com/devesh2997/consequent/errorx"

var (
	errUserAlreadyExistsForEmail = func() error {
		return errorx.NewBusinessError(-1, "user already exists for the given email")
	}
	errInvalidOTP = func() error {
		return errorx.NewBusinessError(-1, "invalid otp")
	}
	errInvalidMobile = func() error {
		return errorx.NewBusinessError(-1, "invalid mobile number")
	}
	errInvalidEmail = func() error {
		return errorx.NewBusinessError(-1, "invalid email")
	}
	errWrongPassword = func() error {
		return errorx.NewBusinessError(-1, "wrong password")
	}
	errOTPHasExpired = func() error {
		return errorx.NewBusinessError(-1, "otp has expired")
	}
	errUserNotFoundForEmail = func() error {
		return errorx.NewBusinessError(-1, "user not found for email")
	}
	errInvalidPassword = func() error {
		return errorx.NewBusinessError(-1, "invalid password")
	}
)
