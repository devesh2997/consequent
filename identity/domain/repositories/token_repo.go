package repositories

import (
	"context"

	"github.com/devesh2997/consequent/identity/domain/entities"
)

type TokenRepo interface {
	SaveRefreshToken(ctx context.Context, token entities.RefreshToken) error
	GetRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	GetPrivateKey() ([]byte, error)
	GetPublicKey() ([]byte, error)
}
