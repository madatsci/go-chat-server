package api

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
