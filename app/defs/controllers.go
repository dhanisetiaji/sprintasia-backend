package defs

import (
	"gotham/controllers"
	"gotham/policies"
	"gotham/services"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"
)

var ControllersDefs = []dingo.Def{
	{
		Name:  "user-controller",
		Scope: di.App,
		Build: func(service services.IUserService, userPolicy policies.IUserPolicy) (controllers.UserController, error) {
			return controllers.UserController{
				UserService: service,
				UserPolicy:  userPolicy,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("user-service"),
			"1": dingo.Service("user-policy"),
		},
	},
	{
		Name:  "auth-controller",
		Scope: di.App,
		Build: func(service services.IAuthService) (controllers.AuthController, error) {
			return controllers.AuthController{
				AuthService: service,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("auth-service"),
		},
	},
	{
		Name:  "task-controller",
		Scope: di.App,
		Build: func(service services.ITaskService) (controllers.TaskController, error) {
			return controllers.TaskController{
				TaskService: service,
			}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("task-service"),
		},
	},
}
