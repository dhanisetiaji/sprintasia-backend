package defs

import (
	"gotham/infrastructures"
	"gotham/repositories"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"
)

var RepositoriesDefs = []dingo.Def{
	{
		Name:  "user-repository",
		Scope: di.App,
		Build: func(gormDatabase infrastructures.IGormDatabase) (repositories.IUserRepository, error) {
			return &repositories.UserRepository{IGormDatabase: gormDatabase}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("db"),
		},
	},
	{
		Name:  "task-repository",
		Scope: di.App,
		Build: func(gormDatabase infrastructures.IGormDatabase) (repositories.IUTaskListRepository, error) {
			return &repositories.TaskListRepository{IGormDatabase: gormDatabase}, nil
		},
		Params: dingo.Params{
			"0": dingo.Service("db"),
		},
	},
}
