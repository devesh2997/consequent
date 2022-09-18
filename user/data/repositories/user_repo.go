package repositories

import (
	"context"
	"errors"

	"github.com/devesh2997/consequent/errorx"
	"github.com/devesh2997/consequent/user/data/mappers"
	"github.com/devesh2997/consequent/user/data/models"
	"github.com/devesh2997/consequent/user/domain/entities"
	"github.com/devesh2997/consequent/user/domain/repositories"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return userRepo{db: db}
}

func (repo userRepo) Create(ctx context.Context, user entities.User) (*entities.User, error) {
	userModel := mappers.NewUserMapper().ToModel(user)
	err := repo.db.Create(&userModel).Error
	if err != nil {
		return nil, err
	}

	userEntity := mappers.NewUserMapper().ToEntity(userModel)

	return &userEntity, nil
}

func (repo userRepo) Update(ctx context.Context, user entities.User) error {
	if user.ID == 0 {
		return errorx.NewSystemError(-1, errors.New("user id is required"))
	}
	err := repo.db.Save(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo userRepo) FindByID(ctx context.Context, id int64) (*entities.User, error) {
	userModel := models.User{}
	err := repo.db.Find(&userModel, id).Error
	if err != nil {
		return nil, err
	}

	userEntity := mappers.NewUserMapper().ToEntity(userModel)

	return &userEntity, nil
}

func (repo userRepo) FindByMobile(ctx context.Context, mobile string) (*entities.User, error) {
	userModel := models.User{}
	err := repo.db.Where("mobile = ?", mobile).Find(&userModel).Error
	if err != nil {
		return nil, err
	}

	userEntity := mappers.NewUserMapper().ToEntity(userModel)

	return &userEntity, nil
}

func (repo userRepo) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	userModel := models.User{}
	err := repo.db.Where("email = ?", email).Find(&userModel).Error
	if err != nil {
		return nil, err
	}

	userEntity := mappers.NewUserMapper().ToEntity(userModel)

	return &userEntity, nil
}
