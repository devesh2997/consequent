package containers

import (
	"github.com/devesh2997/consequent/config"
	"github.com/devesh2997/consequent/datasources"
	"github.com/devesh2997/consequent/identity/data/repositories"
	"github.com/devesh2997/consequent/identity/domain/services"
	"github.com/devesh2997/consequent/identity/presentation/controllers"
	"github.com/devesh2997/consequent/otpsender"
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
	otpSender := otpsender.New2FactorOTPSender(config.Config.Factor2Config.APIKey, config.Config.Factor2Config.OTPTemplateName)

	return services.NewIdentityService(repo, userService, tokenService, otpSender)
}

func InjectIdentityController() controllers.IdentityController {
	return controllers.NewIdentityController(InjectIdentityService())
}
