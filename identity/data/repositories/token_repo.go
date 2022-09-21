package repositories

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/devesh2997/consequent/identity/data/mappers"
	"github.com/devesh2997/consequent/identity/data/models"
	"github.com/devesh2997/consequent/identity/domain/entities"
	"github.com/devesh2997/consequent/identity/domain/repositories"
	"gorm.io/gorm"
)

type tokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) repositories.TokenRepo {
	return tokenRepo{db: db}
}

func (repo tokenRepo) SaveRefreshToken(ctx context.Context, token entities.RefreshToken) error {
	tokenModel := mappers.NewRefreshTokenMapper().ToModel(token)
	tokenModel.UpdatedAt = time.Now()
	err := repo.db.Save(&tokenModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo tokenRepo) GetRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	refreshToken := models.RefreshToken{}
	err := repo.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}

	refreshTokenEntity := mappers.NewRefreshTokenMapper().ToEntity(refreshToken)

	return &refreshTokenEntity, nil
}

func (repo tokenRepo) GetPrivateKey() ([]byte, error) {
	return ioutil.ReadFile("identity/keys/1_private.pem") // TODO (devesh2997) | a better approach needed for this
}

func (repo tokenRepo) GetPublicKey() ([]byte, error) {
	return ioutil.ReadFile("identity/keys/1_public.pem") // TODO (devesh2997) | a better approach needed for this
}
