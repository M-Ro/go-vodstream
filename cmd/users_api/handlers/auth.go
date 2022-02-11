package handlers

import (
	"encoding/json"
	"errors"
	"git.thorn.sh/Thorn/go-vodstream/api"
	"git.thorn.sh/Thorn/go-vodstream/internal/domain"
	"git.thorn.sh/Thorn/go-vodstream/storage/sql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"time"
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
	ErrTokenInvalid            = errors.New("token is not valid")
)

type AuthHandlerConfig struct {
	SigningSecret string
	IssuerIdent   string
	TokenDuration time.Duration
}

type AuthHandler struct {
	config      AuthHandlerConfig
	userStorage *sql.UserStorage
}

func (h *AuthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/v1/auth/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/v1/auth/register", h.Register).Methods(http.MethodPost)
	r.HandleFunc("/v1/auth/reset_password", h.ResetPassword).Methods(http.MethodPost)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	username := "testuser"

	claims := domain.AuthClaim{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.config.TokenDuration)),
			Issuer:    h.config.IssuerIdent,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(h.config.SigningSecret))
	if err != nil {
		log.Errorf("login failed, token signing error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := api.LoginResponse{ // FIXME
		Success:     true,
		Errors:      []string{},
		UserID:      0,
		AccessToken: signedToken,
	}

	encoded, err := json.Marshal(response)
	if err != nil {
		log.Errorf("Login failed %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(encoded)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Register failed: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var registerRequest api.RegisterRequest
	err = json.Unmarshal(body, &registerRequest)
	if err != nil {
		log.Errorf("Register failed: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

}

func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {

}

func getConfig() AuthHandlerConfig {
	viper.SetDefault("api.auth.signing_secret", "change")
	viper.SetDefault("api.auth.ident", "vodstream")
	viper.SetDefault("api.auth.token_duration", "6h")

	return AuthHandlerConfig{
		SigningSecret: viper.GetString("api.auth.signing_secret"),
		IssuerIdent:   viper.GetString("api.auth.ident"),
		TokenDuration: viper.GetDuration("api.auth.token_duration"),
	}
}

func NewAuthHandler(userStorage *sql.UserStorage) AuthHandler {
	return AuthHandler{
		config:      getConfig(),
		userStorage: userStorage,
	}
}
