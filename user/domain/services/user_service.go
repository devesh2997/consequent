package services

import (
	"context"

	"github.com/devesh2997/consequent/user/domain/entities"
	"github.com/devesh2997/consequent/user/domain/repositories"
)

type UserService interface {
	Create(ctx context.Context, user entities.User) (*entities.User, error)
	Update(ctx context.Context, user entities.User) error
	FindByID(ctx context.Context, id int64) (*entities.User, error)
	FindByMobile(ctx context.Context, mobile string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}

func NewUserService(repo repositories.UserRepository) UserService {
	return userService{repo: repo}
}

type userService struct {
	repo repositories.UserRepository
}

func (service userService) Create(ctx context.Context, user entities.User) (*entities.User, error) {
	return service.repo.Create(ctx, user)
}

func (service userService) Update(ctx context.Context, user entities.User) error {
	existingUser, err := service.repo.FindByID(ctx, user.ID)
	if err != nil && err != repositories.ErrUserNotFound {
		return err
	}
	if existingUser == nil {
		return errUserNotFound()
	}

	return service.repo.Update(ctx, user)
}

func (service userService) FindByID(ctx context.Context, id int64) (*entities.User, error) {
	return service.repo.FindByID(ctx, id)
}

func (service userService) FindByMobile(ctx context.Context, mobile string) (*entities.User, error) {
	return service.repo.FindByMobile(ctx, mobile)
}

func (service userService) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	return service.repo.FindByEmail(ctx, email)
}
