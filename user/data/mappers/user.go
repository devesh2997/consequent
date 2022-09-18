package mappers

import (
	"github.com/devesh2997/consequent/user/data/models"
	"github.com/devesh2997/consequent/user/domain/entities"
)

type userMapper struct{}

func NewUserMapper() userMapper {
	return userMapper{}
}

func (mapper userMapper) ToModel(entity entities.User) models.User {
	return models.User{
		ID:     entity.ID,
		Mobile: entity.Mobile,
		Email:  entity.Email,
		Name:   entity.Name,
		Gender: entity.Gender,
	}
}

func (mapper userMapper) ToEntity(model models.User) entities.User {
	return entities.User{
		ID:     model.ID,
		Mobile: model.Mobile,
		Email:  model.Email,
		Name:   model.Name,
		Gender: model.Gender,
	}
}
