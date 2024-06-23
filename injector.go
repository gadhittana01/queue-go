//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gadhittana-01/form-go/app"
	querier "github.com/gadhittana-01/form-go/db/repository"
	"github.com/gadhittana-01/form-go/handler"
	"github.com/gadhittana-01/form-go/service"
	"github.com/gadhittana-01/form-go/utils"
	"github.com/go-chi/chi"
	"github.com/google/wire"
)

var userHandlerSet = wire.NewSet(
	querier.NewRepository,
	handler.NewUserHandler,
	service.NewUserSvc,
)

func InitializeApp(
	route *chi.Mux,
	DB utils.PGXPool,
	config *utils.BaseConfig,
) (app.App, error) {
	wire.Build(
		userHandlerSet,
		app.NewApp,
	)

	return nil, nil
}
