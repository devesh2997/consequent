package services

import (
	"context"
	"crypto/rand"
	"io"
	"strconv"
	"time"

	"github.com/devesh2997/consequent/errorx"
	"github.com/devesh2997/consequent/identity/constants"
	"github.com/devesh2997/consequent/identity/domain/entities"
	"github.com/devesh2997/consequent/identity/domain/repositories"
	userEntities "github.com/devesh2997/consequent/user/domain/entities"
	userRepositories "github.com/devesh2997/consequent/user/domain/repositories"
	"github.com/devesh2997/consequent/user/domain/services"
	"github.com/google/uuid"
)

const (
	otpExpiryDuration = time.Minute * 10
)

type IdentityService interface {
	SendOTP(ctx context.Context, mobileNumber int64) (verificationID string, err error)
	VerifyOTP(ctx context.Context, verificationID string, mobileNumber int64, otp int) (*entities.Token, error)
}

func NewIdentityService(repo repositories.IdentityRepo, userService services.UserService) IdentityService {
	return identityService{repo: repo, userService: userService, tokenGenerator: NewTokenGenerator()}
}

type identityService struct {
	repo           repositories.IdentityRepo
	userService    services.UserService
	otpSender      OTPSender
	tokenGenerator TokenGenerator
}

func (identityService) generateOTP(numDigits int) (int, error) {
	b := make([]byte, numDigits)
	n, err := io.ReadAtLeast(rand.Reader, b, numDigits)
	if n != numDigits {
		return 0, err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	otpStr := string(b)
	otp, err := strconv.Atoi(otpStr)
	if err != nil {
		return 0, err
	}

	return otp, nil
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func (service identityService) SendOTP(ctx context.Context, mobileNumber int64) (verificationID string, err error) {
	otp, err := service.generateOTP(4)
	if err != nil {
		return "", errorx.NewSystemError(-1, err)
	}

	verificationID = uuid.New().String()
	err = service.repo.StoreUserLoginMobileOTP(ctx, entities.UserLoginMobileOTP{
		VerificationID: verificationID,
		Mobile:         mobileNumber,
		OTP:            otp,
		Status:         constants.USER_LOGIN_MOBILE_OTP_STATUS_ACTIVE,
		CreatedAt:      time.Now(),
		ExpiryAt:       time.Now().Add(otpExpiryDuration),
	})
	if err != nil {
		return "", errorx.NewSystemError(-1, err)
	}

	if err := service.otpSender.Send(ctx, mobileNumber, otp); err != nil {
		return "", err
	}

	return verificationID, nil
}

func (service identityService) VerifyOTP(ctx context.Context, verificationID string, mobileNumber int64, otp int) (*entities.Token, error) {
	if err := service.verifyOTP(ctx, verificationID, mobileNumber, otp); err != nil {
		return nil, err
	}

	user, err := service.userService.FindByMobile(ctx, mobileNumber)
	if err != nil && err != userRepositories.ErrUserNotFound {
		return nil, err
	}

	if err == userRepositories.ErrUserNotFound {
		user, err = service.userService.Create(ctx, userEntities.User{
			Mobile: mobileNumber,
		})
		if err != nil {
			return nil, err
		}
	}

	return service.tokenGenerator.Generate(*user)
}

func (service identityService) verifyOTP(ctx context.Context, verificationID string, mobileNumber int64, otp int) error {
	userLoginMobileOTP, err := service.repo.GetUserLoginMobileOTP(ctx, verificationID)
	if err != nil {
		return errorx.NewSystemError(-1, err)
	}

	if userLoginMobileOTP.Mobile != mobileNumber {
		return errorx.NewSystemError(-1, errInvalidMobile())
	}

	if userLoginMobileOTP.OTP != otp {
		return errInvalidOTP()
	}

	nextStatus := constants.USER_LOGIN_MOBILE_OTP_STATUS_VERIFIED
	if time.Now().After(userLoginMobileOTP.ExpiryAt) {
		nextStatus = constants.USER_LOGIN_MOBILE_OTP_STATUS_EXPIRED
		err = errOTPHasExpired()
	}

	userLoginMobileOTP.Status = nextStatus
	if e := service.repo.StoreUserLoginMobileOTP(ctx, *userLoginMobileOTP); e != nil {
		err = errorx.NewSystemError(-1, e)
	}

	return err
}

func (identityService) SignUpWithEmail(ctx context.Context, email string, password string) error {

	return nil
}