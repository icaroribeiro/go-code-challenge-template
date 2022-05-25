package dbtrx_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/model"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/response"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	dbtrxmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/dbtrx"
	domainfactorymodel "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/core/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestMiddlewareUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestNewContext() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrxCtxValue := &gorm.DB{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInCreatingACopyOfAContextWithAnAssociatedValue",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = db
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedCtx := dbtrxmiddlewarepkg.NewContext(ctx, dbTrxCtxValue)

			if !tc.WantError {
				assert.NotEmpty(t, returnedCtx)
				returnedDBTrxCtxValue, ok := dbtrxmiddlewarepkg.FromContext(returnedCtx)
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.Equal(t, dbTrxCtxValue, returnedDBTrxCtxValue)
			}
		})
	}
}

func (ts *TestSuite) TestFromContext() {
	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrxCtxValue := &gorm.DB{}

	ctx := context.Background()

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAnAssociatedValueFromAContext",
			SetUp: func(t *testing.T) {
				dbTrxCtxValue = db
				ctx = dbtrxmiddlewarepkg.NewContext(ctx, dbTrxCtxValue)
			},
			WantError: false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedDBTrxCtxValue, ok := dbtrxmiddlewarepkg.FromContext(ctx)

			if !tc.WantError {
				assert.True(t, ok, "Unexpected type assertion error.")
				assert.NotEmpty(t, returnedDBTrxCtxValue)
				assert.Equal(t, dbTrxCtxValue, returnedDBTrxCtxValue)
			}
		})
	}
}

func (ts *TestSuite) TestDBTrx() {
	user := domainfactorymodel.NewUser(nil)

	body := fmt.Sprintf(`
	{
		"username":"%s",
	}`,
		user.Username)

	driver := "postgres"
	db, mock := NewMockDB(driver)
	dbAux := &gorm.DB{}

	var handlerFunc func(w http.ResponseWriter, r *http.Request)

	statusCode := 0

	sqlQuery := `INSERT INTO "users" ("id","username","created_at","updated_at") VALUES ($1,$2,$3,$4)`

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithDBTrxMiddleware",
			SetUp: func(t *testing.T) {
				dbAux = db

				statusCode = http.StatusOK

				handlerFunc = func(w http.ResponseWriter, r *http.Request) {
					dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

					userDatastore := datastoremodel.User{
						Username: user.Username,
					}

					_ = dbAux.Create(&userDatastore)

					responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: "ok"})
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
			WantError: false,
		},
		{
			Context: "ItShouldFailIfTheDatabaseParameterUsedByTheDBTrxMiddlewareIsNil",
			SetUp: func(t *testing.T) {
				dbAux = nil

				statusCode = http.StatusInternalServerError

				handlerFunc = func(w http.ResponseWriter, r *http.Request) {
					_, ok := dbtrxmiddlewarepkg.FromContext(r.Context())
					if !ok {
						responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
						return
					}

					responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: "ok"})
				}
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionPerformedByTheWrappedFunctionFails",
			SetUp: func(t *testing.T) {
				dbAux = db

				statusCode = http.StatusInternalServerError

				handlerFunc = func(w http.ResponseWriter, r *http.Request) {
					dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

					userDatastore := datastoremodel.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
					}

					responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: "ok"})
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheCommitStatementToEndTheDatabaseTransactionExecutedInsideTheDBTrxMiddlewareFails",
			SetUp: func(t *testing.T) {
				dbAux = db

				statusCode = http.StatusOK

				handlerFunc = func(w http.ResponseWriter, r *http.Request) {
					dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

					userDatastore := datastoremodel.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
					}

					responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: "ok"})
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit().WillReturnError(customerror.New("failed"))
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheRollbackStatementToEndTheDatabaseTransactionExecutedInsideTheDBTrxMiddlewareFails",
			SetUp: func(t *testing.T) {
				dbAux = db

				statusCode = http.StatusInternalServerError

				handlerFunc = func(w http.ResponseWriter, r *http.Request) {
					dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

					userDatastore := datastoremodel.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
					}

					responsehttputilpkg.RespondWithJson(w, http.StatusOK, messagehttputilpkg.Message{Text: "ok"})
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback().WillReturnError(customerror.New("failed"))
			},
			WantError: true,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionPerformedByTheWrappedFunctionFailsAndTheFunctionCallsPanicMethodWithErrorParameterToStopItsExecutionImmediately",
			SetUp: func(t *testing.T) {
				dbAux = db

				statusCode = http.StatusInternalServerError

				handlerFunc = func(w http.ResponseWriter, r *http.Request) {
					dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

					userDatastore := datastoremodel.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						// It is duplicated only to test the code that evaluates
						// if the header is already written in the WriteHeader method.
						responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
						responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
					}

					panic(customerror.New("failed"))
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()
			},
			WantError:   true,
			ShouldPanic: true,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionPerformedByTheWrappedFunctionFailsAndTheFunctionCallsPanicMethodWithNonErrorParameterToStopItsExecutionImmediately",
			SetUp: func(t *testing.T) {
				dbAux = db

				statusCode = http.StatusInternalServerError

				handlerFunc = func(w http.ResponseWriter, r *http.Request) {
					dbAux, _ := dbtrxmiddlewarepkg.FromContext(r.Context())

					userDatastore := datastoremodel.User{
						Username: user.Username,
					}

					result := dbAux.Create(&userDatastore)
					if result.Error != nil {
						// It is duplicated only to test the code that evaluates
						// if the header is already written in the WriteHeader method.
						responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
						responsehttputilpkg.RespondErrorWithJson(w, customerror.New("failed"))
					}

					panic("failed")
				}

				mock.ExpectBegin()

				mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
					WithArgs(sqlmock.AnyArg(), user.Username, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(customerror.New("failed"))

				mock.ExpectRollback()
			},
			WantError:   true,
			ShouldPanic: true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			dbtrxMiddleware := dbtrxmiddlewarepkg.DBTrx(dbAux)

			returnedHandlerFunc := adapterhttputilpkg.AdaptFunc(handlerFunc).With(dbtrxMiddleware)

			route := routehttputilpkg.Route{
				Name:        "Testing",
				Method:      http.MethodGet,
				Path:        "/testing",
				HandlerFunc: returnedHandlerFunc,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
				Body:   body,
			}

			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Result().Header.Get("Content-Type"), "application/json")
				assert.Equal(t, statusCode, resprec.Result().StatusCode)
				returnedMessage := messagehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
				assert.NotEmpty(t, returnedMessage.Text)
			} else {
				if tc.ShouldPanic {
					shouldPanic(t, handlerFunc, resprec, req)
				} else {
					assert.Equal(t, statusCode, resprec.Result().StatusCode)
				}
			}

			err := mock.ExpectationsWereMet()
			assert.Nil(ts.T(), err, fmt.Sprintf("There were unfulfilled expectations: %v.", err))
		})
	}
}

func shouldPanic(t *testing.T, f func(w http.ResponseWriter, r *http.Request), w http.ResponseWriter, r *http.Request) {
	defer func() { recover() }()
	f(w, r)
	t.Errorf("It should have panicked.")
}
