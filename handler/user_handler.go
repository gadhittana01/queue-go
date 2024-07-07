package handler

import (
	"net/http"

	"github.com/gadhittana-01/queue-go/dto"
	"github.com/gadhittana-01/queue-go/service"
	"github.com/gadhittana-01/queue-go/utils"
	"github.com/go-chi/chi"
)

type UserHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)

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

func (h *UserHandlerImpl) SignUp(w http.ResponseWriter, r *http.Request) {
	input := utils.ValidateBodyPayload(r.Body, &dto.SignUpReq{})

	resp := h.userSvc.SignUp(r.Context(), input)

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func (h *UserHandlerImpl) SignIn(w http.ResponseWriter, r *http.Request) {
	input := utils.ValidateBodyPayload(r.Body, &dto.SignInReq{})

	resp := h.userSvc.SignIn(r.Context(), input)

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func setupUserV1Routes(route *chi.Mux, h *UserHandlerImpl) {
	// auth
	route.Group(func(r chi.Router) {
		route.Post("/v1/sign-up", h.SignUp)
		route.Post("/v1/sign-in", h.SignIn)
	})
}
