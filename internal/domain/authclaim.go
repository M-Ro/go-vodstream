package domain

import "github.com/golang-jwt/jwt/v4"

type AuthClaim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
