package containers

import (
	"github.com/devesh2997/consequent/datasources"
	"github.com/devesh2997/consequent/user/data/repositories"
	"github.com/devesh2997/consequent/user/domain/services"
	"github.com/devesh2997/consequent/user/presentation/controllers"
)

func InjectUserService() services.UserService {
	ds, err := datasources.Get()
	if err != nil {
		panic(err)
	}

	repo := repositories.NewUserRepository(ds.SQLClients.GetGormDB())

	return services.NewUserService(repo)
}

func InjectUserController() controllers.UserController {
	service := InjectUserService()

	return controllers.NewUserController(service)
}
