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
	mockauth "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockauth"
	mocksecuritypkg "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mocksecurity"
	mockvalidator "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockvalidator"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestLogOut() {
	id := ""

	errorType := customerror.NoType

	tokenExpTimeInSec := fake.Number(2, 10)

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInLoggingOut",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()

				authID, err := uuid.FromString(id)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				userID := uuid.NewV4()

				args := map[string]interface{}{
					"id":     authID,
					"userID": userID,
				}

				auth := domainentity.AuthFactory(args)

				returnArgs = ReturnArgs{
					{nil},
					{auth, nil},
				}
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheIDIsNotValid",
			SetUp: func(t *testing.T) {
				id = ""

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
					{domainentity.Auth{}, nil},
				}

				errorType = customerror.BadRequest
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenDeletingTheAuth",
			SetUp: func(t *testing.T) {
				id = uuid.NewV4().String()

				returnArgs = ReturnArgs{
					{nil},
					{domainentity.Auth{}, customerror.New("failed")},
				}

				errorType = customerror.NoType
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			validator := new(mockvalidator.Validator)
			validator.On("ValidateWithTags", id, "nonzero, uuid").Return(returnArgs[0]...)

			persistentAuthRepository := new(authdatastoremockrepository.Repository)
			persistentAuthRepository.On("Delete", id).Return(returnArgs[1]...)

			persistentUserRepository := new(userdatastoremockrepository.Repository)

			persistentLoginRepository := new(logindatastoremockrepository.Repository)

			authN := new(mockauth.Auth)
			security := new(mocksecuritypkg.Security)

			authService := authservice.New(persistentAuthRepository, persistentLoginRepository, persistentUserRepository,
				authN, security, validator, tokenExpTimeInSec)

			err := authService.LogOut(id)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
			}
		})
	}
}
