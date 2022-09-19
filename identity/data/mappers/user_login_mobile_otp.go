package mappers

import (
	"github.com/devesh2997/consequent/identity/data/models"
	"github.com/devesh2997/consequent/identity/domain/entities"
)

type userLoginMobileOTP struct{}

func NewUserLoginMobileOTP() userLoginMobileOTP {
	return userLoginMobileOTP{}
}

func (userLoginMobileOTP) ToModel(entity entities.UserLoginMobileOTP) models.UserLoginMobileOTP {
	return models.UserLoginMobileOTP{
		ID:             entity.ID,
		VerificationID: entity.VerificationID,
		Mobile:         entity.Mobile,
		OTP:            entity.OTP,
		Status:         entity.Status,
		CreatedAt:      entity.CreatedAt,
		ExpiryAt:       entity.ExpiryAt,
		UpdatedAt:      entity.UpdatedAt,
	}
}

func (userLoginMobileOTP) ToEntity(model models.UserLoginMobileOTP) entities.UserLoginMobileOTP {
	return entities.UserLoginMobileOTP{
		ID:             model.ID,
		VerificationID: model.VerificationID,
		Mobile:         model.Mobile,
		OTP:            model.OTP,
		Status:         model.Status,
		CreatedAt:      model.CreatedAt,
		ExpiryAt:       model.ExpiryAt,
		UpdatedAt:      model.UpdatedAt,
	}
}
