package services

import (
	"github.com/devesh2997/consequent/identity/domain/entities"
	userEntities "github.com/devesh2997/consequent/user/domain/entities"
)

type TokenGenerator interface {
	Generate(user userEntities.User) (*entities.Token, error)
}

type tokenGenerator struct{}

// Generate implements TokenGenerator
func (tokenGenerator) Generate(user userEntities.User) (*entities.Token, error) {
	panic("unimplemented")
}

func NewTokenGenerator() TokenGenerator {
	return tokenGenerator{}
}
