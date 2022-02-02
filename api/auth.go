package api

// AuthenticationSet represents the auth structured sent along with all authorized requests.
type AuthenticationSet struct {
	AccessToken string `json:"accessToken"`
}

// LoginRequest covers a request to login to a user account from a client.
type LoginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
	RememberMe      bool   `json:"remember"`
}

// LoginResponse covers a response sent from the User API to a client upon login attempt.
type LoginResponse struct {
	Success     bool     `json:"success"`
	Errors      []string `json:"errors"`
	UserID      uint64   `json:"userID"`
	AccessToken string   `json:"token"`
}

// RegisterRequest covers a request to register a user account from a client.
type RegisterRequest struct {
	Email           string `json:"email"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
	AcceptedTerms   bool   `json:"acceptedTerms"`
}

type RegisterStatus int

const (
	RegistrationPendingEmailVerification = iota
	RegistrationPendingAdminVerification
	RegistrationFailed
	RegistrationComplete
)

// RegisterResponse covers a response sent from the User API to a client upon registration.
type RegisterResponse struct {
	Status   RegisterStatus `json:"status"`
	Errors   []string       `json:"errors"`
	Messages []string       `json:"messages"`
}

type ChangePasswordRequest struct {
	Auth            AuthenticationSet `json:"auth"`
	Password        string            `json:"password"`
	ConfirmPassword string            `json:"confirmPassword"`
}

type ChangePasswordResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
}
