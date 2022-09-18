package mappers

import (
	"github.com/devesh2997/consequent/identity/data/models"
	"github.com/devesh2997/consequent/identity/domain/entities"
)

type userPasswordMapper struct{}

func NewUserPasswordMapper() userPasswordMapper {
	return userPasswordMapper{}
}

func (userPasswordMapper) ToModel(entity entities.UserPassword) models.UserPassword {
	return models.UserPassword{
		ID:       entity.ID,
		UserID:   entity.UserID,
		Password: entity.Password,
		Status:   entity.Status,
	}
}

func (userPasswordMapper) ToEntity(model models.UserPassword) entities.UserPassword {
	return entities.UserPassword{
		ID:       model.ID,
		UserID:   model.UserID,
		Password: model.Password,
		Status:   model.Status,
	}
}
