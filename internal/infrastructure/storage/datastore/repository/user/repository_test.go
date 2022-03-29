package user_test

// import (
// 	"errors"
// 	"fmt"
// 	"regexp"
// 	"testing"
// 	"time"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
// 	userdbmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/model/user"
// 	userdsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/postgres/repository/user"
// 	userdbmodelfactory "github.com/icaroribeiro/go-code-challenge-template/tests/factory/infrastructure/persistence/datastore/model/user"
// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/gorm"
// )

// func TestRepository(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestCreate() {
// 	var user userdbmodel.User

// 	var newUser userdbmodel.User

// 	errorType := customerror.NoType

// 	sqlQuery := `INSERT INTO "users" ("id","username","created_at","updated_at") VALUES ($1,$2,$3,$4)`

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInCreatingTheUser",
// 			SetUp: func(t *testing.T) {
// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				user = userdbmodelfactory.New(args)

// 				args = map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"username":  user.Username,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				newUser = userdbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
// 					WillReturnResult(sqlmock.NewResult(1, 1))

// 				ts.SQLMock.ExpectCommit()
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheUser",
// 			SetUp: func(t *testing.T) {
// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				user = userdbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
// 					WillReturnError(errors.New("failed"))

// 				ts.SQLMock.ExpectRollback()

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenCreatingTheUserBecauseTheUserIsAlreadyRegistered",
// 			SetUp: func(t *testing.T) {
// 				args := map[string]interface{}{
// 					"id":        uuid.Nil,
// 					"createdAt": time.Time{},
// 					"updatedAt": time.Time{},
// 				}

// 				user = userdbmodelfactory.New(args)

// 				ts.SQLMock.ExpectBegin()

// 				ts.SQLMock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
// 					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
// 					WillReturnError(customerror.Conflict.New("duplicate key value"))

// 				ts.SQLMock.ExpectRollback()

// 				errorType = customerror.Conflict
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			userDSRepository := userdsrepository.New(ts.DB)

// 			returnedUser, err := userDSRepository.Create(user)

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, newUser.Username, returnedUser.Username)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedUser)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestGetAll() {
// 	var user userdbmodel.User

// 	errorType := customerror.NoType

// 	sqlQuery := `SELECT * FROM "users"`

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingAllUsers",
// 			SetUp: func(t *testing.T) {
// 				user = userdbmodelfactory.New(nil)

// 				rows := sqlmock.
// 					NewRows([]string{"id", "username", "created_at", "updated_at"}).
// 					AddRow(user.ID, user.Username, user.CreatedAt, user.UpdatedAt)

// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WillReturnRows(rows)
// 			},
// 			WantError: false,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenFindingAllUser",
// 			SetUp: func(t *testing.T) {
// 				ts.SQLMock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
// 					WillReturnError(errors.New("failed"))

// 				errorType = customerror.NoType
// 			},
// 			WantError: true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			userDSRepository := userdsrepository.New(ts.DB)

// 			returnedUsers, err := userDSRepository.GetAll()

// 			if !tc.WantError {
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, user.ID, returnedUsers[0].ID)
// 				assert.Equal(t, user.Username, returnedUsers[0].Username)
// 				assert.Equal(t, user.CreatedAt, returnedUsers[0].CreatedAt)
// 				assert.Equal(t, user.UpdatedAt, returnedUsers[0].UpdatedAt)
// 			} else {
// 				assert.NotNil(t, err, "Predicted error lost.")
// 				assert.Equal(t, errorType, customerror.GetType(err))
// 				assert.Empty(t, returnedUsers)
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

// 			userDSRepository := userdsrepository.New(ts.DB)

// 			returnedUserDSRepository := userDSRepository.WithDBTrx(dbTrx)

// 			if !tc.WantError {
// 				assert.NotEmpty(t, returnedUserDSRepository, "Repository interface is empty.")
// 				assert.Equal(t, userDSRepository, returnedUserDSRepository, "Repository interfaces are not the same.")
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) AfterTest(_, _ string) {
// 	err := ts.SQLMock.ExpectationsWereMet()
// 	assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
// }
