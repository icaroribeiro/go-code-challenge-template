package auth_test

// import (
// 	"errors"
// 	"fmt"
// 	"regexp"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/icaroribeiro/new-go-code-challenge-template/internal/application/customerror"
// 	authdbmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/database/model/auth"
// 	authdbrepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/database/repository/auth"
// 	authdbmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/infrastructure/storage/database/model/auth"
// 	logindbmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/infrastructure/storage/database/model/login"
// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/gorm"
// )

// func TestRepository(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestCreate() {
// 	var auth authdbmodel.Auth

// 	var newAuth authdbmodel.Auth

// 	errorType := customerror.NoType

// 	sqlQuery := `INSERT INTO "auths" ("id","user_id","created_at") VALUES ($1,$2,$3)`

// 	sqlQuery2 := `SELECT * FROM "logins" WHERE user_id=$1`

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInCreatingTheAuth",
// 			SetUp: func(t *testing.T) {
// 				userID := uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    userID,
// 					"createdAt": time.Time{},
// 				}

// 				auth = authdbmodelfactory.New(args)

// 				args = map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    auth.UserID,
// 					"createdAt": time.Time{},
// 				}

// 				newAuth = authdbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
// 					WillReturnResult(sqlmock.NewResult(1, 1))

// 				ts.SQLMock.ExpectCommit()
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheAuth",
// 			SetUp: func(t *testing.T) {
// 				userID := uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    userID,
// 					"createdAt": time.Time{},
// 				}

// 				auth = authdbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
// 					WillReturnError(errors.New("failed"))

// 				ts.SQLMock.ExpectRollback()

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheAuthBecauseTheUserAuthIsAlreadyRegistered",
// 			SetUp: func(t *testing.T) {
// 				userID := uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    userID,
// 					"createdAt": time.Time{},
// 				}

// 				auth = authdbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
// 					WillReturnError(customerror.Conflict.New("auths_user_id_key"))

// 				ts.SQLMock.ExpectRollback()

// 				args = map[string]interface{}{
// 					"userID": userID,
// 				}
// 				login := logindbmodelfactory.New(args)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
// 					AddRow(login.ID, login.UserID, login.Username, login.Password, login.CreatedAt, login.UpdatedAt)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery2)).
// 					WithArgs(auth.UserID).
// 					WillReturnRows(rows)

// 				errorType = customerror.Conflict
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheAuthBecauseTheUserAuthIsAlreadyRegisteredAndLoginIsNotFound",
// 			SetUp: func(t *testing.T) {
// 				userID := uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    userID,
// 					"createdAt": time.Time{},
// 				}

// 				auth = authdbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
// 					WillReturnError(errors.New("auths_user_id_key"))

// 				ts.SQLMock.ExpectRollback()

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery2)).
// 					WithArgs(auth.UserID).
// 					WillReturnRows(&sqlmock.Rows{})

// 				errorType = customerror.NotFound
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheAuthBecauseTheUserAuthIsAlreadyRegisteredAndAnErrorAlsoHappensWhenFindingTheLoginByUserID",
// 			SetUp: func(t *testing.T) {
// 				userID := uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    userID,
// 					"createdAt": time.Time{},
// 				}

// 				auth = authdbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), auth.UserID, sqlmock.AnyArg()).
// 					WillReturnError(errors.New("auths_user_id_key"))

// 				ts.SQLMock.ExpectRollback()

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery2)).
// 					WithArgs(auth.UserID).
// 					WillReturnError(errors.New("failed"))

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDBRepository := authdbrepository.New(ts.DB)

// 			returnedAuth, err := authDBRepository.Create(auth)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, newAuth.UserID, returnedAuth.UserID)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedAuth)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestDelete() {
// 	var id uuid.UUID

// 	var auth authdbmodel.Auth

// 	errorType := customerror.NoType

// 	sqlQuery := `SELECT * FROM "auths" WHERE id=$1`

// 	sqlQuery2 := `DELETE FROM "auths" WHERE "auths"."id" = $1`

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInDeletingTheAuth",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id": id,
// 				}

// 				auth = authdbmodelfactory.New(args)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "created_at"}).
// 					AddRow(auth.ID, auth.UserID, auth.CreatedAt)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(id).
// 					WillReturnRows(rows)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery2)).
// 					WithArgs(id).
// 					WillReturnResult(sqlmock.NewResult(0, 1))

// 				ts.SQLMock.ExpectCommit()
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenFindingTheAuthByID",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(id).
// 					WillReturnError(errors.New("failed"))

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthIsNotFound",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(id).
// 					WillReturnRows(&sqlmock.Rows{})

// 				errorType = customerror.NotFound
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenDeletingTheAuthByID",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id": id,
// 				}

// 				auth = authdbmodelfactory.New(args)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "created_at"}).
// 					AddRow(auth.ID, auth.UserID, auth.CreatedAt)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(id).
// 					WillReturnRows(rows)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery2)).
// 					WithArgs(id).
// 					WillReturnError(errors.New("failed"))

// 				ts.SQLMock.ExpectRollback()

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthIsNotDeleted",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id": id,
// 				}

// 				auth = authdbmodelfactory.New(args)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "created_at"}).
// 					AddRow(auth.ID, auth.UserID, auth.CreatedAt)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(id).
// 					WillReturnRows(rows)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery2)).
// 					WithArgs(id).
// 					WillReturnResult(sqlmock.NewResult(0, 0))

// 				ts.SQLMock.ExpectCommit()

// 				errorType = customerror.NotFound
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDBRepository := authdbrepository.New(ts.DB)

// 			returnedAuth, err := authDBRepository.Delete(id.String())

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, auth.ID, returnedAuth.ID)
// 				assert.Equal(t, auth.UserID, returnedAuth.UserID)
// 				assert.Equal(t, auth.CreatedAt, returnedAuth.CreatedAt)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedAuth)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestWithDBTrx() {
// 	dbTrx := &gorm.DB{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSettingTheRepositoryWithDatabaseTransaction",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldSucceedInSettingTheRepositoryWithoutDatabaseTransaction",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = nil
// 			},
// 			WantError: false,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDBRepository := authdbrepository.New(ts.DB)

// 			returnedAuthDBRepository := authDBRepository.WithDBTrx(dbTrx)

// 			if !tc.WantError {
// 				assert.NotEmpty(t, returnedAuthDBRepository, "Repository interface is empty.")
// 				assert.Equal(t, authDBRepository, returnedAuthDBRepository, "Repository interfaces are not the same.")
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) AfterTest(_, _ string) {
// 	err := ts.SQLMock.ExpectationsWereMet()
// 	assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
// }
