package mappers

import (
	"github.com/devesh2997/consequent/identity/data/models"
	"github.com/devesh2997/consequent/identity/domain/entities"
)

type refreshTokenMapper struct{}

func NewRefreshTokenMapper() refreshTokenMapper {
	return refreshTokenMapper{}
}

func (refreshTokenMapper) ToEntity(model models.RefreshToken) entities.RefreshToken {
	return entities.RefreshToken{
		ID:        model.ID,
		Token:     model.Token,
		Status:    model.Status,
		CreatedAt: model.CreatedAt,
		ExpiryAt:  model.ExpiryAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (refreshTokenMapper) ToModel(entity entities.RefreshToken) models.RefreshToken {
	return models.RefreshToken{
		ID:        entity.ID,
		Token:     entity.Token,
		Status:    entity.Status,
		CreatedAt: entity.CreatedAt,
		ExpiryAt:  entity.ExpiryAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

type tokenMapper struct{}

func NewTokenMapper() tokenMapper {
	return tokenMapper{}
}

func (tokenMapper) ToEntity(model models.Token) entities.Token {
	return entities.Token{
		JWT: entities.JWT{
			Token:    model.JWT.Token,
			ExpiryAt: model.JWT.ExpiryAt,
		},
		RefreshToken: NewRefreshTokenMapper().ToEntity(model.RefreshToken),
	}
}

func (tokenMapper) ToModel(entity entities.Token) models.Token {
	return models.Token{
		JWT: models.JWT{
			Token:    entity.JWT.Token,
			ExpiryAt: entity.JWT.ExpiryAt,
		},
		RefreshToken: NewRefreshTokenMapper().ToModel(entity.RefreshToken),
	}
}
