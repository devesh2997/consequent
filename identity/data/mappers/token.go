package mappers

import (
	"github.com/devesh2997/consequent/identity/data/models"
	"github.com/devesh2997/consequent/identity/domain/entities"
)

type tokenMapper struct{}

func NewTokenMapper() tokenMapper {
	return tokenMapper{}
}

func (mapper tokenMapper) ToEntity(model models.Token) entities.Token {
	return entities.Token{
		JWT:          model.JWT,
		RefreshToken: model.RefreshToken,
	}
}

func (mapper tokenMapper) ToModel(entity entities.Token) models.Token {
	return models.Token{
		JWT:          entity.JWT,
		RefreshToken: entity.RefreshToken,
	}
}
