package auth_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// import (
// 	"errors"
// 	"fmt"
// 	"testing"
// 	"time"

// 	fake "github.com/brianvoe/gofakeit/v5"
// 	"github.com/dgrijalva/jwt-go"
// 	"github.com/icaroribeiro/new-go-code-challenge-template/internal/application/customerror"
// 	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model/auth"
// 	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/auth"
// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

func TestAuth(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestCreateToken() {
	rsaKeys := ts.RSAKeys

	auth := domainmodel.Auth{}

	tokenExpTimeInSec := 0

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingAToken",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				tokenExpTimeInSec = fake.Number(30, 60)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authpkg := authpkg.New(rsaKeys)

			tokenString, err := authpkg.CreateToken(auth, tokenExpTimeInSec)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.NotEmpty(t, tokenString, "Unexpected empty token")
			}
		})
	}
}

// func (ts *TestSuite) TestVerifyToken() {
// 	auth := domainmodel.Auth{}

// 	issuedAt := time.Now().Unix()

// 	expiredAt := time.Now().Unix()

// 	rsaKeys := ts.RSAKeys

// 	authpkg := authpkg.New(rsaKeys)

// 	err := errors.New("")

// 	tokenString := ""

// 	isToRefreshToken := false

// 	timeBeforeTokenExpTimeInSec := 60

// 	errorType := customerror.NoType

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInVerifyingAToken",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = domainmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				issuedAt = time.Now().Unix()
// 				tokenExpTimeInSec := fake.Number(30, 60)
// 				duration := time.Second * time.Duration(tokenExpTimeInSec)
// 				expiredAt = time.Now().Add(duration).Unix()

// 				tokenString, err = authpkg.CreateToken(auth, tokenExpTimeInSec)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
// 				assert.NotEmpty(t, tokenString, "")
// 			},
// 			WantError: false,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheTokenHasNoSigningMethodAllowed",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = domainmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				issuedAt = time.Now().Unix()
// 				duration := time.Second * time.Duration(fake.Number(30, 60))
// 				expiredAt = time.Now().Add(duration).Unix()

// 				claims := jwt.MapClaims{
// 					"authID":     auth.ID,
// 					"userID":     auth.UserID,
// 					"iat":        issuedAt,
// 					"exp":        expiredAt,
// 					"authorized": true,
// 				}

// 				token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
// 				tokenString, err = token.SignedString(jwt.UnsafeAllowNoneSignatureType)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
// 				assert.NotEmpty(t, tokenString, "")

// 				isToRefreshToken = false

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheTokenHasExpiredAndItWillNotBeRefreshed",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = domainmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				tokenExpTimeInSec := fake.Number(-10, -2)

// 				tokenString, err = authpkg.CreateToken(auth, tokenExpTimeInSec)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
// 				assert.NotEmpty(t, tokenString, "")

// 				isToRefreshToken = false

// 				errorType = customerror.BadRequest
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheTokenIsInvalid",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = domainmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				tokenString = fake.Word()

// 				isToRefreshToken = true

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheTokenWouldBeRefreshedAndItsExpirationTimeIsAnImproperlyFormattedFloatValue",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = domainmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				issuedAt = time.Now().Unix()
// 				expiredAt := fake.Date()

// 				claims := jwt.MapClaims{
// 					"authID":     auth.ID,
// 					"userID":     auth.UserID,
// 					"iat":        issuedAt,
// 					"exp":        expiredAt,
// 					"authorized": true,
// 				}

// 				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 				tokenString, err = token.SignedString(ts.RSAKeys.PrivateKey)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
// 				assert.NotEmpty(t, tokenString, "")

// 				isToRefreshToken = true

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheTokenHasExpiredAndItWouldBeRefreshedButItsExpirationTimeIsNotWithinTheTimePriorToTheExpirationTime",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = domainmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				tokenExpTimeInSec := fake.Number(90, 120)

// 				tokenString, err = authpkg.CreateToken(auth, tokenExpTimeInSec)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
// 				assert.NotEmpty(t, tokenString, "")

// 				isToRefreshToken = true

// 				errorType = customerror.BadRequest
// 			},
// 			WantError: true,
// 			TearDown:  func(t *testing.T) {},
// 		},
// 	}
// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			token, err := authpkg.VerifyToken(tokenString, isToRefreshToken, timeBeforeTokenExpTimeInSec)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
// 				assert.NotEmpty(t, tokenString, "Unexpected empty token")
// 				claims, ok := token.Claims.(jwt.MapClaims)
// 				assert.True(t, ok, "Unexpected type assertion error")
// 				assert.Equal(t, auth.ID.String(), claims["auth_id"])
// 				assert.Equal(t, auth.UserID.String(), claims["user_id"])
// 				iat, ok := claims["iat"].(float64)
// 				assert.True(t, ok, "Unexpected type assertion error")
// 				assert.WithinDuration(t, time.Unix(issuedAt, 0), time.Unix(int64(iat), 0), time.Second)
// 				if !isToRefreshToken {
// 					exp, ok := claims["exp"].(float64)
// 					assert.True(t, ok, "Unexpected type assertion error")
// 					assert.WithinDuration(t, time.Unix(expiredAt, 0), time.Unix(int64(exp), 0), time.Second)
// 				}
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Nil(t, token, "Token is not nil")
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestFetchAuth() {
// 	rsaKeys := ts.RSAKeys
// 	authpkg := authpkg.New(rsaKeys)

// 	id := uuid.NewV4()
// 	userID := uuid.NewV4()

// 	token := &jwt.Token{}

// 	errorType := customerror.NoType

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInFetchingAuthDataFromAToken",
// 			SetUp: func(t *testing.T) {
// 				claims := jwt.MapClaims{
// 					"auth_id": id.String(),
// 					"user_id": userID.String(),
// 				}

// 				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 				assert.NotNil(t, token, "Token is nil")
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthIDValueFromTokenIsNotAString",
// 			SetUp: func(t *testing.T) {
// 				claims := jwt.MapClaims{
// 					"auth_id": fake.Number(1, 10),
// 				}

// 				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 				assert.NotNil(t, token, "Token is nil")

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthIDValueFromTokenIsNotAUUIDString",
// 			SetUp: func(t *testing.T) {
// 				claims := jwt.MapClaims{
// 					"auth_id": fake.Word(),
// 				}

// 				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 				assert.NotNil(t, token, "Token is nil")

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheUserIDValueFromTokenIsNotAString",
// 			SetUp: func(t *testing.T) {
// 				claims := jwt.MapClaims{
// 					"auth_id": id.String(),
// 					"user_id": fake.Number(1, 10),
// 				}

// 				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 				assert.NotNil(t, token, "Token is nil")

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheUserIDValueFromTokenIsNotAUUIDString",
// 			SetUp: func(t *testing.T) {
// 				claims := jwt.MapClaims{
// 					"auth_id": id.String(),
// 					"user_id": fake.Word(),
// 				}

// 				token = jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
// 				assert.NotNil(t, token, "Token is nil")

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 	}
// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			auth, err := authpkg.FetchAuth(token)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
// 				assert.Equal(t, id, auth.ID)
// 				assert.Equal(t, userID, auth.UserID)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, auth, "Auth is not empty")
// 			}
// 		})
// 	}
// }
