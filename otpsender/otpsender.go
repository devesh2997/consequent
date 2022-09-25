package otpsender

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devesh2997/consequent/errorx"
)

type OTPSender interface {
	Send(ctx context.Context, mobileNumber string, otp int) error
}

func New2FactorOTPSender(apiKey string, otpTemplateName string) OTPSender {
	return factor2{
		apiKey:          apiKey,
		otpTemplateName: otpTemplateName,
	}
}

type factor2 struct {
	apiKey          string
	otpTemplateName string
}

func (f2 factor2) Send(ctx context.Context, mobileNumber string, otp int) error {
	url := f2.getSendOTPURL(mobileNumber, otp)
	res, err := http.Get(url) // TODO (devesh2997) | better handling of the response object can be done here
	if err != nil {
		return errorx.NewSystemError(-1, err)
	}
	fmt.Println(res.StatusCode)
	// fmt.Printf("sending otp: %d\n", otp)

	return nil
}

func (f2 factor2) getSendOTPURL(mobileNumber string, otp int) string {
	baseURL := "https://2factor.in/API/V1/"

	return fmt.Sprintf("%s/%s/SMS/%s/%d/%s", baseURL, f2.apiKey, mobileNumber, otp, f2.otpTemplateName)
}
