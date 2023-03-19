package user_test

import (
	"testing"

	userservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/user"
	userdatastoremockrepository "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/infrastructure/datastore/mockrepository/user"
	"github.com/icaroribeiro/go-code-challenge-template/tests/mocks/pkg/mockvalidator"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestWithDBTrx() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrx := &gorm.DB{}

	persistentUserRepositoryWithDBTrx := &userdatastoremockrepository.Repository{}

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSettingRepositoriesWithDatabaseTransaction",
			SetUp: func(t *testing.T) {
				dbTrx = db

				persistentUserRepositoryWithDBTrx = &userdatastoremockrepository.Repository{}

				returnArgs = ReturnArgs{
					{persistentUserRepositoryWithDBTrx},
				}
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			persistentUserRepository := new(userdatastoremockrepository.Repository)
			persistentUserRepository.On("WithDBTrx", dbTrx).Return(returnArgs[0]...)

			validator := new(mockvalidator.Validator)

			userService := userservice.New(persistentUserRepository, validator)

			returnedUserService := userService.WithDBTrx(dbTrx)

			if !tc.WantError {
				assert.NotEmpty(t, returnedUserService, "Service interface is empty.")
				assert.Equal(t, userService, returnedUserService, "Service interfaces are not the same.")
			}
		})
	}
}
