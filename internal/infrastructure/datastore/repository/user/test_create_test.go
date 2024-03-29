package user_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domainentity "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/entity"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/user"
	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestCreate() {
	driver := "postgres"
	db, mock := NewMockDB(driver)

	user := domainentity.User{}

	newUser := domainentity.User{}

	errorType := customerror.NoType

	sqlQuery := `INSERT INTO "users" ("id","username","created_at","updated_at") VALUES ($1,$2,$3,$4)`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingTheUser",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				user = domainentity.UserFactory(args)

				args = map[string]interface{}{
					"id":       uuid.Nil,
					"username": user.Username,
				}

				newUser = domainentity.UserFactory(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheUser",
			SetUp: func(t *testing.T) {
				args := map[string]interface{}{
					"id": uuid.Nil,
				}

				user = domainentity.UserFactory(args)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()

				errorType = customerror.NoType
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheUserBecauseTheUserIsAlreadyRegistered",
			SetUp: func(t *testing.T) {
				user = domainentity.UserFactory(nil)

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.Conflict.New("duplicate key value"))

				mock.ExpectRollback()

				errorType = customerror.Conflict
			},
			WantError: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			persistentUserRepository := userdatastorerepository.New(db)

			returnedUser, err := persistentUserRepository.Create(user)

			if !tc.WantError {
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, newUser.Username, returnedUser.Username)
			} else {
				assert.NotNil(t, err, "Predicted error lost.")
				assert.Equal(t, errorType, customerror.GetType(err))
				assert.Empty(t, returnedUser)
			}

			err = mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}
