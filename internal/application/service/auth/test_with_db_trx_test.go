package auth_test

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/auth"
	authdatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/auth"
	logindatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/login"
	userdatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/user"
	mockauth "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockauth"
	mocksecuritypkg "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mocksecurity"
	mockvalidator "github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockvalidator"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestWithDBTrx() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrx := &gorm.DB{}

	persistentAuthRepositoryWithDBTrx := &authdatastoremockrepository.Repository{}
	persistentUserRepositoryWithDBTrx := &userdatastoremockrepository.Repository{}
	persistentLoginRepositoryWithDBTrx := &logindatastoremockrepository.Repository{}

	tokenExpTimeInSec := fake.Number(2, 10)

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSettingRepositoriesWithDatabaseTransaction",
			SetUp: func(t *testing.T) {
				dbTrx = db

				persistentAuthRepositoryWithDBTrx = &authdatastoremockrepository.Repository{}
				persistentUserRepositoryWithDBTrx = &userdatastoremockrepository.Repository{}
				persistentLoginRepositoryWithDBTrx = &logindatastoremockrepository.Repository{}

				returnArgs = ReturnArgs{
					{persistentAuthRepositoryWithDBTrx},
					{persistentUserRepositoryWithDBTrx},
					{persistentLoginRepositoryWithDBTrx},
				}
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			persistentAuthRepository := new(authdatastoremockrepository.Repository)
			persistentAuthRepository.On("WithDBTrx", dbTrx).Return(returnArgs[0]...)

			persistentUserRepository := new(userdatastoremockrepository.Repository)
			persistentUserRepository.On("WithDBTrx", dbTrx).Return(returnArgs[1]...)

			persistentLoginRepository := new(logindatastoremockrepository.Repository)
			persistentLoginRepository.On("WithDBTrx", dbTrx).Return(returnArgs[2]...)

			authN := new(mockauth.Auth)
			security := new(mocksecuritypkg.Security)
			validator := new(mockvalidator.Validator)

			authService := authservice.New(persistentAuthRepository, persistentLoginRepository, persistentUserRepository,
				authN, security, validator, tokenExpTimeInSec)

			returnedAuthService := authService.WithDBTrx(dbTrx)

			if !tc.WantError {
				assert.NotEmpty(t, returnedAuthService, "Service interface is empty.")
				assert.Equal(t, authService, returnedAuthService, "Service interfaces are not the same.")
			}
		})
	}
}
