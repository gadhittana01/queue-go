package app

import (
	"fmt"
	"net/http"

	"github.com/gadhittana-01/queue-go/handler"
	"github.com/gadhittana-01/queue-go/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type App interface {
	Start()
}

type AppImpl struct {
	route        *chi.Mux
	config       *utils.BaseConfig
	userHandler  handler.UserHandler
	queueHandler handler.QueueHandler
}

func NewApp(route *chi.Mux,
	config *utils.BaseConfig,
	userHandler handler.UserHandler,
	queueHandler handler.QueueHandler,
) App {
	return &AppImpl{
		route:        route,
		config:       config,
		userHandler:  userHandler,
		queueHandler: queueHandler,
	}
}

func (s *AppImpl) Start() {
	s.route.Use(utils.Recovery)
	s.route.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	s.userHandler.SetupUserRoutes(s.route)
	s.queueHandler.SetupQueueRoutes(s.route)

	s.route.NotFound(func(w http.ResponseWriter, r *http.Request) {
		utils.GenerateErrorResp[any](w, nil, 404)
	})

	utils.LogInfo(fmt.Sprintf("server started on port %d", s.config.ServerPort))
	port := fmt.Sprintf(":%d", s.config.ServerPort)
	err := http.ListenAndServe(port, s.route)
	if err != nil {
		panic(err)
	}
}
