package defs

import (
	"gotham/repositories"
	"gotham/services"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"
)

var ServicesDefs = []dingo.Def{
	{
		Name:  "auth-service",
		Scope: di.App,
		Build: func(repository repositories.IUserRepository) (s services.IAuthService, err error) {
			s = &services.AuthService{UserRepository: repository}
			return s, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("user-repository"),
		},
	},
	{
		Name:  "user-service",
		Scope: di.App,
		Build: func(repository repositories.IUserRepository) (s services.IUserService, err error) {
			return &services.UserService{UserRepository: repository}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("user-repository"),
		},
	},
	{
		Name:  "task-service",
		Scope: di.App,
		Build: func(repository repositories.IUTaskListRepository) (s services.ITaskService, err error) {
			return &services.TaskService{TaskListRepository: repository}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("task-repository"),
		},
	},
}
