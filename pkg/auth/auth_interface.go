package auth

import (
	"github.com/golang-jwt/jwt"
	domainentity "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
)

// IAuth interface is the auth's contract.
type IAuth interface {
	CreateToken(auth domainentity.Auth, tokenExpTimeInSec int) (string, error)
	ExtractTokenString(authHeaderString string) (string, error)
	DecodeToken(tokenString string) (*jwt.Token, error)
	ValidateTokenRenewal(token *jwt.Token, timeBeforeTokenExpTimeInSec int) (*jwt.Token, error)
	FetchAuthFromToken(token *jwt.Token) (domainentity.Auth, error)
}
