package services

import (
	"context"
	"crypto/rand"
	"io"
	"net/mail"
	"strconv"
	"time"

	"github.com/devesh2997/consequent/errorx"
	"github.com/devesh2997/consequent/identity/constants"
	"github.com/devesh2997/consequent/identity/domain/entities"
	"github.com/devesh2997/consequent/identity/domain/repositories"
	"github.com/devesh2997/consequent/logger"
	"github.com/devesh2997/consequent/otpsender"
	userEntities "github.com/devesh2997/consequent/user/domain/entities"
	userRepositories "github.com/devesh2997/consequent/user/domain/repositories"
	"github.com/devesh2997/consequent/user/domain/services"
	"github.com/google/uuid"
)

const (
	otpExpiryDuration = time.Minute * 10
)

type IdentityService interface {
	SendOTP(ctx context.Context, mobileNumber string) (verificationID string, err error)
	VerifyOTP(ctx context.Context, verificationID string, mobileNumber string, otp int) (*entities.Token, error)
	ResendOTP(ctx context.Context, verificationID string) (string, error)
	IsEmailRegistered(ctx context.Context, email string) (bool, error)
	SignUpWithEmail(ctx context.Context, email string, password string) (*entities.Token, error)
	SignInWithEmailAndPassword(ctx context.Context, email string, password string) (*entities.Token, error)
}

func NewIdentityService(repo repositories.IdentityRepo, userService services.UserService, tokenService TokenService, otpSender otpsender.OTPSender) IdentityService {
	return identityService{repo: repo, userService: userService, tokenService: tokenService, otpSender: otpSender}
}

type identityService struct {
	repo         repositories.IdentityRepo
	userService  services.UserService
	otpSender    otpsender.OTPSender
	tokenService TokenService
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

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

func (service identityService) SendOTP(ctx context.Context, mobileNumber string) (verificationID string, err error) {
	if err := service.validateMobile(mobileNumber); err != nil {
		return "", err
	}
	otp, err := service.generateOTP(4)
	if err != nil {
		return "", errorx.NewSystemError(-1, err)
	}

	verificationID = uuid.New().String()
	err = service.repo.SaveUserLoginMobileOTP(ctx, entities.UserLoginMobileOTP{
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

	go service.sendOTP(mobileNumber, otp)

	return verificationID, nil
}

func (service identityService) sendOTP(mobileNumber string, otp int) {
	ctx := context.TODO()
	if err := service.otpSender.Send(ctx, mobileNumber, otp); err != nil {
		logger.Log.Error(ctx, err)
	}
}

func (service identityService) validateMobile(mobileNumber string) error {
	if len(mobileNumber) != 10 { // TODO (devesh2997) | this validation can be improved
		return errInvalidMobile()
	}

	return nil
}

func (service identityService) VerifyOTP(ctx context.Context, verificationID string, mobileNumber string, otp int) (*entities.Token, error) {
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

	return service.tokenService.Generate(ctx, *user)
}

func (service identityService) ResendOTP(ctx context.Context, verificationID string) (string, error) {
	userLoginMobileOTP, err := service.repo.GetUserLoginMobileOTP(ctx, verificationID)
	if err != nil {
		return "", errorx.NewSystemError(-1, err)
	}
	if !userLoginMobileOTP.IsActive() || userLoginMobileOTP.HasExpired() { // TODO (devesh2997) | mark the old otp as expired if neccessary
		return service.SendOTP(ctx, userLoginMobileOTP.Mobile)
	}

	go service.sendOTP(userLoginMobileOTP.Mobile, userLoginMobileOTP.OTP)

	return verificationID, nil
}

func (service identityService) verifyOTP(ctx context.Context, verificationID string, mobileNumber string, otp int) error {
	userLoginMobileOTP, err := service.repo.GetUserLoginMobileOTP(ctx, verificationID)
	if err != nil {
		return errorx.NewSystemError(-1, err)
	}

	if userLoginMobileOTP.Mobile != mobileNumber {
		return errInvalidMobile()
	}

	if userLoginMobileOTP.OTP != otp || !userLoginMobileOTP.IsActive() {
		return errInvalidOTP()
	}

	nextStatus := constants.USER_LOGIN_MOBILE_OTP_STATUS_VERIFIED
	if userLoginMobileOTP.HasExpired() {
		nextStatus = constants.USER_LOGIN_MOBILE_OTP_STATUS_EXPIRED
		err = errOTPHasExpired()
	}

	userLoginMobileOTP.Status = nextStatus
	if e := service.repo.SaveUserLoginMobileOTP(ctx, *userLoginMobileOTP); e != nil {
		err = errorx.NewSystemError(-1, e)
	}

	return err
}

func (service identityService) SignUpWithEmail(ctx context.Context, email string, password string) (*entities.Token, error) {
	if err := service.validateEmailAndPassword(email, password); err != nil {
		return nil, err
	}
	existingUser, err := service.userService.FindByEmail(ctx, email)
	if err != nil && err != userRepositories.ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		errUserAlreadyExistsForEmail()
	}

	user, err := service.userService.Create(ctx, userEntities.User{Email: email})
	if err != nil {
		return nil, err
	}

	err = service.repo.SaveUserPassword(ctx, entities.UserPassword{
		UserID:   user.ID,
		Password: password, // TODO (devesh2997) | hashed value should be stored here, not plaintext
		Status:   constants.USER_PASSWORD_STATUS_ACTIVE,
	})
	if err != nil {
		return nil, err
	}

	return service.tokenService.Generate(ctx, *user)
}

func (service identityService) IsEmailRegistered(ctx context.Context, email string) (bool, error) {
	if !service.isEmailValid(email) {
		return false, errInvalidEmail()
	}

	existingUser, err := service.userService.FindByEmail(ctx, email)
	if err != nil && err != userRepositories.ErrUserNotFound {
		return false, err
	}
	if existingUser != nil {
		return true, nil
	}

	return false, nil
}

func (service identityService) SignInWithEmailAndPassword(ctx context.Context, email string, password string) (*entities.Token, error) {
	if err := service.validateEmailAndPassword(email, password); err != nil {
		return nil, err
	}
	existingUser, err := service.userService.FindByEmail(ctx, email)
	if err != nil && err != userRepositories.ErrUserNotFound {
		return nil, err
	}
	if existingUser == nil {
		return nil, errUserNotFoundForEmail()
	}

	userPassword, err := service.repo.GetActiveUserPassword(ctx, existingUser.ID)
	if err != nil {
		return nil, err
	}
	if userPassword.Password != password {
		return nil, errWrongPassword()
	}

	return service.tokenService.Generate(ctx, *existingUser)
}

func (service identityService) validateEmailAndPassword(email string, password string) error {
	if !service.isEmailValid(email) {
		return errInvalidEmail()
	}
	if !service.isPasswordValid(password) {
		return errInvalidPassword()
	}

	return nil
}

func (service identityService) isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)

	return err == nil
}

func (service identityService) isPasswordValid(password string) bool {
	return len(password) > 5 // TODO (devesh2997) | password rules can be much stricter
}
