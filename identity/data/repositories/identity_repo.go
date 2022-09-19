package repositories

import (
	"context"
	"time"

	"github.com/devesh2997/consequent/identity/constants"
	"github.com/devesh2997/consequent/identity/data/mappers"
	"github.com/devesh2997/consequent/identity/data/models"
	"github.com/devesh2997/consequent/identity/domain/entities"
	"github.com/devesh2997/consequent/identity/domain/repositories"
	"gorm.io/gorm"
)

func NewIdentityRepo(db *gorm.DB) repositories.IdentityRepo {
	return identityRepo{db: db}
}

type identityRepo struct {
	db *gorm.DB
}

func (repo identityRepo) GetActiveUserPassword(ctx context.Context, userID int64) (*entities.UserPassword, error) {
	userPassword := models.UserPassword{}
	err := repo.db.Where("status = ?", constants.USER_PASSWORD_STATUS_ACTIVE).First(&userPassword).Error
	if err != nil {
		return nil, err
	}

	userPasswordEntity := mappers.NewUserPasswordMapper().ToEntity(userPassword)

	return &userPasswordEntity, nil
}

func (repo identityRepo) StoreUserPassword(ctx context.Context, userPassword entities.UserPassword) error {
	userPasswordModel := mappers.NewUserPasswordMapper().ToModel(userPassword)
	err := repo.db.Create(&userPasswordModel).Error
	if err != nil {
		return err
	}

	return err
}

func (repo identityRepo) StoreUserLoginMobileOTP(ctx context.Context, otp entities.UserLoginMobileOTP) error {
	model := mappers.NewUserLoginMobileOTP().ToModel(otp)
	model.UpdatedAt = time.Now()
	if err := repo.db.Create(&model).Error; err != nil {
		return err
	}

	return nil
}

func (repo identityRepo) GetUserLoginMobileOTP(ctx context.Context, verificationID string) (*entities.UserLoginMobileOTP, error) {
	otp := models.UserLoginMobileOTP{}
	err := repo.db.Where("verification_id = ?", verificationID).Find(&otp).Error
	if err != nil {
		return nil, err
	}

	entity := mappers.NewUserLoginMobileOTP().ToEntity(otp)

	return &entity, nil
}
