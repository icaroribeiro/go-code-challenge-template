package auth

import (
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	uuid "github.com/satori/go.uuid"
)

type Auth struct {
	RSAKeys RSAKeys
}

// RSAKeys is the representation of the RSA keys.
type RSAKeys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

func New(rsaKeys RSAKeys) IAuth {
	return &Auth{
		RSAKeys: RSAKeys{
			PublicKey:  rsaKeys.PublicKey,
			PrivateKey: rsaKeys.PrivateKey,
		},
	}
}

// CreateToken is the function that creates a new token for a specific auth and time duration.
func (a *Auth) CreateToken(auth domainmodel.Auth, tokenExpTimeInSec int) (string, error) {
	duration := time.Second * time.Duration(tokenExpTimeInSec)

	claims := jwt.MapClaims{
		"auth_id":    auth.ID.String(),
		"user_id":    auth.UserID.String(),
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(duration).Unix(),
		"authorized": true,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(a.RSAKeys.PrivateKey)
}

// VerifyToken is the function that translates the token string in jwt token and checks if the jwt token is valid or not.
func (a *Auth) VerifyToken(tokenString string, isToRefreshToken bool, timeBeforeTokenExpTimeInSec int) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, customerror.New("unexpected signing method when trying to decode the token")
		}
		return a.RSAKeys.PublicKey, nil
	}

	token, err := jwt.Parse(tokenString, keyFunc)

	if verr, ok := err.(*jwt.ValidationError); ok {
		switch verr.Errors {
		case jwt.ValidationErrorExpired:
			if !isToRefreshToken {
				errorMessage := "the token has expired, then perform a refresh action to get a new one"
				return nil, customerror.BadRequest.New(errorMessage)
			}
		default:
			return nil, err
		}
	}

	// claims, ok := token.Claims.(jwt.MapClaims)
	// if !ok {
	// 	errorMessage := "failed to fetch data from the token"
	// 	return nil, customerror.New(errorMessage)
	// }

	// if isToRefreshToken {
	// 	// It is necessary to ensure that a new token will not be issued until enough time has elapsed.
	// 	expiredAt, ok := claims["exp"].(float64)
	// 	if !ok {
	// 		errorMessage := "failed to fetch data from the token"
	// 		return nil, customerror.New(errorMessage)
	// 	}

	// 	duration := time.Second * time.Duration(timeBeforeTokenExpTimeInSec)

	// 	if time.Until(time.Unix(int64(expiredAt), 0)) > duration {
	// 		errorMessage := "the token expiration time is not within the time prior to the expiration time"
	// 		return nil, customerror.BadRequest.New(errorMessage)
	// 	}
	// }

	return token, nil
}

// FetchAuth is the function that get auth data from the token.
func (a *Auth) FetchAuth(token *jwt.Token) (domainmodel.Auth, error) {
	auth := domainmodel.Auth{}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return auth, customerror.New("failed to fetch data from the token")
	}

	id, ok := claims["auth_id"].(string)
	if !ok {
		return auth, customerror.New("failed to extract the auth_id from the token")
	}

	authID, err := uuid.FromString(id)
	if err != nil {
		return auth, customerror.Newf("failed to convert the auth_id %s from the token to UUID", id)
	}

	id, ok = claims["user_id"].(string)
	if !ok {
		return auth, customerror.New("failed to extract the user_id from the token")
	}

	userID, err := uuid.FromString(id)
	if err != nil {
		return auth, customerror.Newf("failed to convert the user_id %s from the to UUID", id)
	}

	auth.ID = authID
	auth.UserID = userID

	return auth, nil
}
