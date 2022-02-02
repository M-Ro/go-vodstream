package handlers

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	ErrLoginInvalidCredentials = errors.New("invalid login credentials provided")
	ErrRegisterMissingDetails  = errors.New("missing details during registration")
	ErrRegisterUsernameExists  = errors.New("user with this username already exists")
	ErrRegisterEmailExists     = errors.New("user with this email already exists")
	ErrRegisterNoTermsAccepted = errors.New("terms and conditions field not accepted")
	ErrPasswordRequirements    = errors.New("password does not meet minimum requirements")
	ErrPasswordMatch           = errors.New("password fields do not match")
	ErrPasswordPrevious        = errors.New("new password matches existing password")
)

type AuthHandler struct {
}

func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/v1/auth/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/v1/auth/register", h.Register).Methods(http.MethodPost)
	r.HandleFunc("/v1/auth/reset_password", h.ResetPassword).Methods(http.MethodPost)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

}
