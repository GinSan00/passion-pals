package model

type UserResponse struct {
	Status    string       `json:"status"`
	Responder *UserProfile `json:"responder"`
}
