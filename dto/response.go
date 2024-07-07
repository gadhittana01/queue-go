package dto

type SignUpRes struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	ExpToken int64  `json:"expToken"`
}

type SignInRes struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Token    string `json:"token"`
	ExpToken int64  `json:"expToken"`
}

type GetQueueRes struct {
	ID               string `json:"id"`
	QueueNumber      string `json:"queueNumber"`
	UserID           string `json:"userId"`
	ArrivalTime      string `json:"arrivalTime"`
	ServiceTime      string `json:"serviceTime"`
	TotalWaitingTime string `json:"totalWaitingTime"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type CreateQueueRes struct {
	ID          string `json:"id"`
	QueueNumber string `json:"queueNumber"`
	UserID      string `json:"userId"`
	ArrivalTime string `json:"arrivalTime"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ProcessQueueRes struct {
	QueueNumber      string `json:"queueNumber"`
	UserID           string `json:"userId"`
	ArrivalTime      string `json:"arrivalTime"`
	ServiceTime      string `json:"serviceTime"`
	TotalWaitingTime string `json:"totalWaitingTime"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type SimpleMessageRes struct {
	Message string `json:"message"`
}
