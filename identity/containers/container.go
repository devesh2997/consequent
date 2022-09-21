package containers

import (
	"github.com/devesh2997/consequent/datasources"
	"github.com/devesh2997/consequent/identity/data/repositories"
	"github.com/devesh2997/consequent/identity/domain/services"
	"github.com/devesh2997/consequent/user/containers"
)

func InjectTokenService() services.TokenService {
	ds, err := datasources.Get()
	if err != nil {
		panic(err)
	}

	repo := repositories.NewTokenRepo(ds.SQLClients.GetGormDB())

	return services.NewTokenService(repo)
}

func InjectIdentityService() services.IdentityService {
	ds, err := datasources.Get()
	if err != nil {
		panic(err)
	}

	repo := repositories.NewIdentityRepo(ds.SQLClients.GetGormDB())
	userService := containers.InjectUserService()
	tokenService := InjectTokenService()

	return services.NewIdentityService(repo, userService, tokenService)
}