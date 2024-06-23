package dto

type CreateUserReq struct {
	Name           string `json:"name" validate:"required"`
	IdentityNumber string `json:"identityNumber" validate:"required"`
	Email          string `json:"email" validate:"required"`
	DateOfBirth    string `json:"dateOfBirth" validate:"required"`
}
