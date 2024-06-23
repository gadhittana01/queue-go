package handler

import (
	"net/http"

	"github.com/gadhittana-01/form-go/dto"
	"github.com/gadhittana-01/form-go/service"
	"github.com/gadhittana-01/form-go/utils"
	"github.com/go-chi/chi"
)

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)

	SetupUserRoutes(route *chi.Mux)
}

type UserHandlerImpl struct {
	userSvc service.UserSvc
}

func NewUserHandler(
	userSvc service.UserSvc,
) UserHandler {
	return &UserHandlerImpl{
		userSvc: userSvc,
	}
}

func (h *UserHandlerImpl) SetupUserRoutes(route *chi.Mux) {
	setupUserV1Routes(route, h)
}

func (h *UserHandlerImpl) CreateUser(w http.ResponseWriter, r *http.Request) {
	input := utils.ValidateBodyPayload(r.Body, &dto.CreateUserReq{})

	resp := h.userSvc.CreateUser(r.Context(), input)

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func setupUserV1Routes(route *chi.Mux, h *UserHandlerImpl) {
	route.Post("/v1/users", h.CreateUser)
}
