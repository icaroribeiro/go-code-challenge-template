package auth

import (
	"github.com/dgrijalva/jwt-go"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
)

// IAuth interface is a collection of function signatures that represents the auth's contract.
type IAuth interface {
	CreateToken(auth domainmodel.Auth, tokenExpTimeInSec int) (string, error)
	VerifyToken(tokenString string, isToRefreshToken bool, timeBeforeTokenExpTimeInSec int) (*jwt.Token, error)
	FetchAuth(token *jwt.Token) (domainmodel.Auth, error)
}
