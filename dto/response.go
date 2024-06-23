package dto

type CreateUserRes struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	IdentityNumber string `json:"identityNumber"`
	Email          string `json:"email"`
	DateOfBirth    string `json:"dateOfBirth"`
}
