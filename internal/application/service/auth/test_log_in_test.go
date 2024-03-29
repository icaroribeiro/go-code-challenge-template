package auth_test

import (
	"fmt"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/auth"
	domainentity "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	authdatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/auth"
	logindatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/login"
	userdatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/user"
	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template/pkg/security"
	mockauth "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockauth"
	mocksecuritypkg "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mocksecurity"
	mockvalidator "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockvalidator"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestLogIn() {
	credentials := securitypkg.Credentials{}

	login := domainentity.Login{}

	auth := domainentity.Auth{}

	newAuth := domainentity.Auth{}

	tokenExpTimeInSec := fake.Number(2, 10)

	token := ""

	errorType := customerror.NoType

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInLoggingIn",
			SetUp: func(t *testing.T) {
				credentials = securitypkg.CredentialsFactory(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainentity.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainentity.Auth{
					UserID: login.UserID,
				}

				id = uuid.NewV4()

				newAuth = domainentity.Auth{
					ID:     id,
					UserID: userID,
				}

				token = fake.Word()

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{domainentity.Auth{}, nil},
					{newAuth, nil},
					{token, nil},
				}
			},
			WantError: false,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheCredentialsAreNotValid",
			SetUp: func(t *testing.T) {
				credentials = securitypkg.Credentials{}

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
					{domainentity.Login{}, nil},
					{nil},
					{domainentity.Auth{}, nil},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingALoginByUsername",
			SetUp: func(t *testing.T) {
				credentials = securitypkg.CredentialsFactory(nil)

				returnArgs = ReturnArgs{
					{nil},
					{domainentity.Login{}, customerror.New("failed")},
					{nil},
					{domainentity.Auth{}, nil},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheUsernameIsNotRegistered",
			SetUp: func(t *testing.T) {
				credentials = securitypkg.CredentialsFactory(nil)

				returnArgs = ReturnArgs{
					{nil},
					{domainentity.Login{}, nil},
					{nil},
					{domainentity.Auth{}, nil},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NotFound
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenVerifyingThePasswords",
			SetUp: func(t *testing.T) {
				credentials = securitypkg.CredentialsFactory(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainentity.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{customerror.New("failed")},
					{domainentity.Auth{}, nil},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAnAuth",
			SetUp: func(t *testing.T) {
				credentials = securitypkg.CredentialsFactory(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainentity.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainentity.Auth{
					UserID: login.UserID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{domainentity.Auth{}, customerror.New("failed")},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheUserIDIsAlreadyRegistered",
			SetUp: func(t *testing.T) {
				credentials = securitypkg.CredentialsFactory(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainentity.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainentity.Auth{
					UserID: login.UserID,
				}

				id = uuid.NewV4()

				newAuth = domainentity.Auth{
					ID:     id,
					UserID: login.UserID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{auth, nil},
					{domainentity.Auth{}, nil},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingANewAuth",
			SetUp: func(t *testing.T) {
				credentials = securitypkg.CredentialsFactory(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainentity.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainentity.Auth{
					UserID: login.UserID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{domainentity.Auth{}, nil},
					{domainentity.Auth{}, customerror.New("failed")},
					{"", nil},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingAToken",
			SetUp: func(t *testing.T) {
				credentials = securitypkg.CredentialsFactory(nil)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				login = domainentity.Login{
					ID:       id,
					UserID:   userID,
					Username: credentials.Username,
					Password: credentials.Password,
				}

				auth = domainentity.Auth{
					UserID: login.UserID,
				}

				id = uuid.NewV4()

				newAuth = domainentity.Auth{
					ID:     id,
					UserID: login.UserID,
				}

				returnArgs = ReturnArgs{
					{nil},
					{login, nil},
					{nil},
					{domainentity.Auth{}, nil},
					{newAuth, nil},
					{"", customerror.New("failed")},
				}

				errorType = customerror.NoType
			},
			WantError: true,
			TearDown:  func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			validator := new(mockvalidator.Validator)
			validator.On("Validate", credentials).Return(returnArgs[0]...)

			persistentLoginRepository := new(logindatastoremockrepository.Repository)
			persistentLoginRepository.On("GetByUsername", credentials.Username).Return(returnArgs[1]...)

			security := new(mocksecuritypkg.Security)
			security.On("VerifyPasswords", login.Password, credentials.Password).Return(returnArgs[2]...)

			persistentAuthRepository := new(authdatastoremockrepository.Repository)
			persistentAuthRepository.On("GetByUserID", login.UserID.String()).Return(returnArgs[3]...)
			persistentAuthRepository.On("Create", auth).Return(returnArgs[4]...)

			authN := new(mockauth.Auth)
			authN.On("CreateToken", newAuth, tokenExpTimeInSec).Return(returnArgs[5]...)

			persistentUserRepository := new(userdatastoremockrepository.Repository)

			authService := authservice.New(persistentAuthRepository, persistentLoginRepository, persistentUserRepository,
				authN, security, validator, tokenExpTimeInSec)

			returnedToken, err := authService.LogIn(credentials)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, token, returnedToken)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedToken)
			}

			tc.TearDown(t)
		})
	}
}
