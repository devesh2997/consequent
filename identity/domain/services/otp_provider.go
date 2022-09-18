package services

import "context"

type OTPSender interface {
	Send(ctx context.Context, mobileNumber int64, otp int) error
}
