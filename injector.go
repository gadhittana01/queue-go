//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gadhittana-01/queue-go/app"
	querier "github.com/gadhittana-01/queue-go/db/repository"
	"github.com/gadhittana-01/queue-go/handler"
	"github.com/gadhittana-01/queue-go/service"
	"github.com/gadhittana-01/queue-go/utils"
	"github.com/go-chi/chi"
	"github.com/google/wire"
)

var userHandlerSet = wire.NewSet(
	querier.NewRepository,
	utils.NewToken,
	handler.NewUserHandler,
	service.NewUserSvc,
)

var queueHandlerSet = wire.NewSet(
	handler.NewQueueHandler,
	service.NewQueueSvc,
)

var authMiddlewareSet = wire.NewSet(
	utils.NewAuthMiddleware,
)

func InitializeApp(
	route *chi.Mux,
	DB utils.PGXPool,
	config *utils.BaseConfig,
) (app.App, error) {
	wire.Build(
		userHandlerSet,
		queueHandlerSet,
		authMiddlewareSet,
		app.NewApp,
	)

	return nil, nil
}
