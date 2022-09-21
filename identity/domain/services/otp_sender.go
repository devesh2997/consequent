package services

import (
	"context"
	"fmt"
)

type OTPSender interface {
	Send(ctx context.Context, mobileNumber string, otp int) error
}

func NewOTPSender() OTPSender {
	return otpSender{}
}

type otpSender struct{}

func (otpSender) Send(ctx context.Context, mobileNumber string, otp int) error {
	fmt.Println("sending otp")

	return nil
}
