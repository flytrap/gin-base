//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
)

func BuildInjector() (*Injector, func(), error) {
	wire.Build(
		InitGormDB,
		InitStore,
		repositories.RepositorySet,
		InitGinEngine,
		services.ServiceSet,
		api.APISet,
		router.RouterSet,
		InjectorSet,
	)
	return new(Injector), nil, nil
}
