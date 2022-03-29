package login_test

// import (
// 	"errors"
// 	"fmt"
// 	"regexp"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	fake "github.com/brianvoe/gofakeit/v5"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/application/customerror"
// 	logindbmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/database/model/login"
// 	logindbrepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/database/repository/login"
// 	securitypkg "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/security"
// 	logindbmodelfactory "github.com/icaroribeiro/go-code-challenge-template/tests/factory/infrastructure/persistence/database/model/login"
// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/gorm"
// )

// func TestRepository(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestCreate() {
// 	var login logindbmodel.Login

// 	var newLogin logindbmodel.Login

// 	errorType := customerror.NoType

// 	sqlQuery := `INSERT INTO "logins" ("id","user_id","username","password","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6)`

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInCreatingTheLogin",
// 			SetUp: func(t *testing.T) {
// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				args = map[string]interface{}{
// 					"userID":   login.UserID,
// 					"username": login.Username,
// 					"password": login.Password,
// 				}
// 				newLogin = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
// 					WillReturnResult(sqlmock.NewResult(1, 1))

// 				ts.SQLMock.ExpectCommit()
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheLogin",
// 			SetUp: func(t *testing.T) {
// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
// 					WillReturnError(errors.New("failed"))

// 				ts.SQLMock.ExpectRollback()

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheLoginBecauseTheUserLoginIsAlreadyRegistered",
// 			SetUp: func(t *testing.T) {
// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
// 					WillReturnError(customerror.Conflict.New("logins_user_id_key"))

// 				ts.SQLMock.ExpectRollback()

// 				errorType = customerror.Conflict
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			loginDSRepository := logindbrepository.New(ts.DB)

// 			returnedLogin, err := loginDSRepository.Create(login)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, newLogin.UserID, returnedLogin.UserID)
// 				assert.Equal(t, newLogin.Username, returnedLogin.Username)
// 				security := securitypkg.New()
// 				err := security.VerifyPasswords(returnedLogin.Password, newLogin.Password)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedLogin)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestGetByUsername() {
// 	var username string

// 	var login logindbmodel.Login

// 	errorType := customerror.NoType

// 	sqlQuery := `SELECT * FROM "logins" WHERE username=$1`

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingTheLoginByUsername",
// 			SetUp: func(t *testing.T) {
// 				username = fake.Username()

// 				args := map[string]interface{}{
// 					"username": username,
// 				}

// 				login = logindbmodelfactory.New(args)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
// 					AddRow(login.ID, login.UserID, login.Username, login.Password, login.CreatedAt, login.UpdatedAt)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(username).
// 					WillReturnRows(rows)
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldSucceedIfTheUsernameIsNotFound",
// 			SetUp: func(t *testing.T) {
// 				username = fake.Username()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    uuid.Nil,
// 					"username":  "",
// 					"password":  "",
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(username).
// 					WillReturnRows(&sqlmock.Rows{})
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenFindingTheLoginByUsername",
// 			SetUp: func(t *testing.T) {
// 				username = fake.Username()

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(username).
// 					WillReturnError(errors.New("failed"))

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			loginDSRepository := logindbrepository.New(ts.DB)

// 			returnedLogin, err := loginDSRepository.GetByUsername(username)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, login.ID, returnedLogin.ID)
// 				assert.Equal(t, login.UserID, returnedLogin.UserID)
// 				assert.Equal(t, login.Username, returnedLogin.Username)
// 				assert.Equal(t, login.Password, returnedLogin.Password)
// 				assert.Equal(t, login.CreatedAt, returnedLogin.CreatedAt)
// 				assert.Equal(t, login.UpdatedAt, returnedLogin.UpdatedAt)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedLogin)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestGetByUserID() {
// 	var userID uuid.UUID

// 	var login logindbmodel.Login

// 	errorType := customerror.NoType

// 	sqlQuery := `SELECT * FROM "logins" WHERE user_id=$1`

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingTheLoginByUserID",
// 			SetUp: func(t *testing.T) {
// 				userID = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"userID": userID,
// 				}
// 				login = logindbmodelfactory.New(args)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
// 					AddRow(login.ID, login.UserID, login.Username, login.Password, login.CreatedAt, login.UpdatedAt)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(userID).
// 					WillReturnRows(rows)
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldSucceedIfTheLoginIsNotFound",
// 			SetUp: func(t *testing.T) {
// 				userID = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"userID":    uuid.Nil,
// 					"username":  "",
// 					"password":  "",
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(userID).
// 					WillReturnRows(&sqlmock.Rows{})
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenFindingTheLoginByUserID",
// 			SetUp: func(t *testing.T) {
// 				userID = uuid.NewV4()

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(userID).
// 					WillReturnError(errors.New("failed"))

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			loginDSRepository := logindbrepository.New(ts.DB)

// 			returnedLogin, err := loginDSRepository.GetByUserID(userID.String())

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, login.ID, returnedLogin.ID)
// 				assert.Equal(t, login.UserID, returnedLogin.UserID)
// 				assert.Equal(t, login.Username, returnedLogin.Username)
// 				assert.Equal(t, login.Password, returnedLogin.Password)
// 				assert.Equal(t, login.CreatedAt, returnedLogin.CreatedAt)
// 				assert.Equal(t, login.UpdatedAt, returnedLogin.UpdatedAt)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedLogin)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestUpdate() {
// 	var id uuid.UUID

// 	var login logindbmodel.Login

// 	var updatedLogin logindbmodel.Login

// 	errorType := customerror.NoType

// 	sqlQuery := `UPDATE "logins" SET "user_id"=$1,"username"=$2,"password"=$3,"updated_at"=$4 WHERE id=$5`

// 	sqlQuery2 := `SELECT * FROM "logins" WHERE id=$1`

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInUpdatingTheLogin",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				args = map[string]interface{}{
// 					"id":       id,
// 					"userID":   login.UserID,
// 					"username": login.Username,
// 					"password": login.Password,
// 				}
// 				updatedLogin = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
// 					WillReturnResult(sqlmock.NewResult(0, 1))

// 				ts.SQLMock.ExpectCommit()

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
// 					AddRow(id, updatedLogin.UserID, updatedLogin.Username, updatedLogin.Password, updatedLogin.CreatedAt, updatedLogin.UpdatedAt)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery2)).
// 					WithArgs(id).
// 					WillReturnRows(rows)
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenUpdatingTheLogin",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
// 					WillReturnError(errors.New("failed"))

// 				ts.SQLMock.ExpectRollback()

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheLoginIsNotFound",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
// 					WillReturnResult(sqlmock.NewResult(0, 0))

// 				ts.SQLMock.ExpectCommit()

// 				errorType = customerror.NotFound
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenFindingTheLoginByID",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
// 					WillReturnResult(sqlmock.NewResult(0, 1))

// 				ts.SQLMock.ExpectCommit()

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery2)).
// 					WithArgs(id).
// 					WillReturnError(errors.New("failed"))

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheLoginIsNotFoundAfterUpdatingIt",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				login = logindbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(login.UserID, login.Username, sqlmock.AnyArg(), sqlmock.AnyArg(), id).
// 					WillReturnResult(sqlmock.NewResult(0, 1))

// 				ts.SQLMock.ExpectCommit()

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery2)).
// 					WithArgs(id).
// 					WillReturnRows(&sqlmock.Rows{})

// 				errorType = customerror.NotFound
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			loginDSRepository := logindbrepository.New(ts.DB)

// 			returnedLogin, err := loginDSRepository.Update(id.String(), login)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, updatedLogin.ID, returnedLogin.ID)
// 				assert.Equal(t, updatedLogin.UserID, returnedLogin.UserID)
// 				assert.Equal(t, updatedLogin.Username, returnedLogin.Username)
// 				assert.Equal(t, updatedLogin.Password, returnedLogin.Password)
// 				assert.Equal(t, updatedLogin.CreatedAt, returnedLogin.CreatedAt)
// 				assert.Equal(t, updatedLogin.UpdatedAt, returnedLogin.UpdatedAt)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedLogin)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestDelete() {
// 	var id uuid.UUID

// 	var login logindbmodel.Login

// 	errorType := customerror.NoType

// 	sqlQuery := `SELECT * FROM "logins" WHERE id=$1`

// 	sqlQuery2 := `DELETE FROM "logins" WHERE "logins"."id" = $1`

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInDeletingTheLogin",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id": id,
// 				}

// 				login = logindbmodelfactory.New(args)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
// 					AddRow(login.ID, login.UserID, login.Username, login.Password, login.CreatedAt, login.UpdatedAt)

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
// 			Context: "ItShouldFailIfAnErrorOccursWhenFindingTheLoginByID",
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
// 			Context: "ItShouldFailIfTheLoginIsNotFound",
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
// 			Context: "ItShouldFailIfAnErrorOccursWhenDeletingTheLoginByID",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id": id,
// 				}
// 				login = logindbmodelfactory.New(args)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
// 					AddRow(login.ID, login.UserID, login.Username, login.Password, login.CreatedAt, login.UpdatedAt)

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
// 			Context: "ItShouldFailIfTheLoginIsNotDeleted",
// 			SetUp: func(t *testing.T) {
// 				id = uuid.NewV4()

// 				args := map[string]interface{}{
// 					"id": id,
// 				}
// 				login = logindbmodelfactory.New(args)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "user_id", "username", "password", "created_at", "updated_at"}).
// 					AddRow(login.ID, login.UserID, login.Username, login.Password, login.CreatedAt, login.UpdatedAt)

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

// 			loginDSRepository := logindbrepository.New(ts.DB)

// 			returnedLogin, err := loginDSRepository.Delete(id.String())

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, login.ID, returnedLogin.ID)
// 				assert.Equal(t, login.UserID, returnedLogin.UserID)
// 				assert.Equal(t, login.Username, returnedLogin.Username)
// 				assert.Equal(t, login.Password, returnedLogin.Password)
// 				assert.Equal(t, login.CreatedAt, returnedLogin.CreatedAt)
// 				assert.Equal(t, login.UpdatedAt, returnedLogin.UpdatedAt)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedLogin)
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

// 			loginDSRepository := logindbrepository.New(ts.DB)

// 			returnedLoginDSRepository := loginDSRepository.WithDBTrx(dbTrx)

// 			if !tc.WantError {
// 				assert.NotEmpty(t, returnedLoginDSRepository, "Repository interface is empty.")
// 				assert.Equal(t, loginDSRepository, returnedLoginDSRepository, "Repository interfaces are not the same.")
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) AfterTest(_, _ string) {
// 	err := ts.SQLMock.ExpectationsWereMet()
// 	assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
// }
