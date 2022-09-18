package repositories

import (
	"context"

	"github.com/devesh2997/consequent/user/domain/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user entities.User) (*entities.User, error)
	Update(ctx context.Context, user entities.User) error
	FindByID(ctx context.Context, id int64) (*entities.User, error)
	FindByMobile(ctx context.Context, mobile string) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
}
