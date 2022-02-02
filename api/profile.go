package api

type ProfileRequest struct {
	Auth   AuthenticationSet `json:"auth"`
	UserID uint64            `json:"userID"`
}

type ProfileResponse struct {
	UserID   uint64 `json:"userID"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
