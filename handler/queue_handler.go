package handler

import (
	"net/http"

	"github.com/gadhittana-01/queue-go/dto"
	"github.com/gadhittana-01/queue-go/service"
	"github.com/gadhittana-01/queue-go/utils"
	"github.com/go-chi/chi"
)

type QueueHandler interface {
	GetQueue(w http.ResponseWriter, r *http.Request)

	SetupQueueRoutes(route *chi.Mux)
}

type QueueHandlerImpl struct {
	queueSvc       service.QueueSvc
	authMiddleware utils.AuthMiddleware
}

func NewQueueHandler(
	queueSvc service.QueueSvc,
	authMiddleware utils.AuthMiddleware,
) QueueHandler {
	return &QueueHandlerImpl{
		queueSvc:       queueSvc,
		authMiddleware: authMiddleware,
	}
}

func (h *QueueHandlerImpl) SetupQueueRoutes(route *chi.Mux) {
	setupQueueV1Routes(route, h)
}

func (h *QueueHandlerImpl) GetQueue(w http.ResponseWriter, r *http.Request) {
	resp := h.queueSvc.GetQueue(r.Context())

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func (h *QueueHandlerImpl) CreateQueue(w http.ResponseWriter, r *http.Request) {
	resp := h.queueSvc.CreateQueue(r.Context())

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func (h *QueueHandlerImpl) ProcessQueue(w http.ResponseWriter, r *http.Request) {
	queueID := utils.ValidateURLParamUUID(r, "queueID")

	resp := h.queueSvc.ProcessQueue(r.Context(), dto.ProcessQueueReq{
		QueueID: queueID,
	})

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func (h *QueueHandlerImpl) DeleteQueue(w http.ResponseWriter, r *http.Request) {
	queueID := utils.ValidateURLParamUUID(r, "queueID")

	resp := h.queueSvc.DeleteQueue(r.Context(), dto.DeleteQueueReq{
		QueueID: queueID,
	})

	utils.GenerateSuccessResp(w, resp, http.StatusOK)
}

func setupQueueV1Routes(route *chi.Mux, h *QueueHandlerImpl) {
	// queue
	route.Group(func(r chi.Router) {
		route.Get("/v1/queue", h.authMiddleware.CheckIsAuthenticated(h.GetQueue))
		route.Post("/v1/queue", h.authMiddleware.CheckIsAuthenticated(h.CreateQueue))
		route.Post("/v1/queue/{queueID}", h.authMiddleware.CheckIsAuthenticated(h.ProcessQueue))
		route.Delete("/v1/queue/{queueID}", h.authMiddleware.CheckIsAuthenticated(h.DeleteQueue))
	})
}
