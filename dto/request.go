package dto

import "github.com/google/uuid"

type SignUpReq struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ProcessQueueReq struct {
	QueueID uuid.UUID `json:"queueId" validate:"required"`
}

type DeleteQueueReq struct {
	QueueID uuid.UUID `json:"queueId" validate:"required"`
}
