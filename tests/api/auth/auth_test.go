package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/application/service/auth"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/model"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/user"
	authhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/handler/auth"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	tokenhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/token"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/security"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestAuthInt(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestSignUp() {
	dbTrx := &gorm.DB{}

	var authN authpkg.IAuth

	credentials := securitypkg.Credentials{}

	body := ""

	contextMap := make(map[interface{}]interface{})

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningUp",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
			SetUp: func(t *testing.T) {
				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = ""

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = nil
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
					"username":"%s",
					"password":"%s"
				`,
					credentials.Username, credentials.Password)

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx
			},
			StatusCode: http.StatusBadRequest,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(dbTrx)

			userDatastoreRepository := userdatastorerepository.New(dbTrx)

			loginDatastoreRepository := logindatastorerepository.New(dbTrx)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "SignUp",
				Method:      http.MethodPost,
				Path:        "/sign_up",
				HandlerFunc: authHandler.SignUp,
			}

			requestData := requesthttputilpkg.RequestData{
				Method:     route.Method,
				Target:     route.Path,
				Body:       body,
				ContextMap: contextMap,
			}

			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

			requesthttputilpkg.SetRequestContext(req, requestData.ContextMap)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedToken := tokenhttputilpkg.Token{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedToken)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
				assert.NotEmpty(t, returnedToken.Text)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}

func (ts *TestSuite) TestSignIn() {
	dbTrx := &gorm.DB{}

	var authN authpkg.IAuth

	userDatastore := datastoremodel.User{}

	loginDatastore := datastoremodel.Login{}

	credentials := securitypkg.Credentials{}

	body := ""

	contextMap := make(map[interface{}]interface{})

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningIn",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				userDatastore = datastoremodel.User{
					Username: username,
				}

				result := dbTrx.Create(&userDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				loginDatastore = datastoremodel.Login{
					UserID:   userDatastore.ID,
					Username: username,
					Password: password,
				}

				result = dbTrx.Create(&loginDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
			SetUp: func(t *testing.T) {
				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = ""

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = nil
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
					"username":"%s",
					"password":"%s"
				`,
					credentials.Username, credentials.Password)

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx
			},
			StatusCode: http.StatusBadRequest,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(dbTrx)

			userDatastoreRepository := userdatastorerepository.New(dbTrx)

			loginDatastoreRepository := logindatastorerepository.New(dbTrx)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "SignIn",
				Method:      http.MethodPost,
				Path:        "/sign_in",
				HandlerFunc: authHandler.SignIn,
			}

			requestData := requesthttputilpkg.RequestData{
				Method:     route.Method,
				Target:     route.Path,
				Body:       body,
				ContextMap: contextMap,
			}

			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

			requesthttputilpkg.SetRequestContext(req, requestData.ContextMap)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedToken := tokenhttputilpkg.Token{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedToken)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
				assert.NotEmpty(t, returnedToken.Text)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}

func (ts *TestSuite) TestRefreshToken() {
	dbTrx := &gorm.DB{}

	var authN authpkg.IAuth

	userDatastore := datastoremodel.User{}

	loginDatastore := datastoremodel.Login{}

	authDatastore := datastoremodel.Auth{}

	contextMap := make(map[interface{}]interface{})

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRefreshingTheToken",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				userDatastore = datastoremodel.User{
					Username: username,
				}

				result := dbTrx.Create(&userDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				loginDatastore = datastoremodel.Login{
					UserID:   userDatastore.ID,
					Username: username,
					Password: password,
				}

				result = dbTrx.Create(&loginDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				authDatastore = datastoremodel.Auth{
					UserID: userDatastore.ID,
				}

				result = dbTrx.Create(&authDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = authDatastore.ToDomain()
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsInvalid",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = fake.Word()
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(dbTrx)

			userDatastoreRepository := userdatastorerepository.New(dbTrx)

			loginDatastoreRepository := logindatastorerepository.New(dbTrx)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "RefreshToken",
				Method:      http.MethodPost,
				Path:        "/refresh_token",
				HandlerFunc: authHandler.RefreshToken,
			}

			requestData := requesthttputilpkg.RequestData{
				Method:     route.Method,
				Target:     route.Path,
				ContextMap: contextMap,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			requesthttputilpkg.SetRequestContext(req, requestData.ContextMap)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedToken := tokenhttputilpkg.Token{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedToken)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
				assert.NotEmpty(t, returnedToken.Text)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}

func (ts *TestSuite) TestChangePassword() {
	dbTrx := &gorm.DB{}

	var authN authpkg.IAuth

	userDatastore := datastoremodel.User{}

	loginDatastore := datastoremodel.Login{}

	authDatastore := datastoremodel.Auth{}

	auth := domainmodel.Auth{}

	passwords := securitypkg.Passwords{}

	body := ""

	contextMap := make(map[interface{}]interface{})

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInResettingThePassword",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				userDatastore = datastoremodel.User{
					Username: username,
				}

				result := dbTrx.Create(&userDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				loginDatastore = datastoremodel.Login{
					UserID:   userDatastore.ID,
					Username: username,
					Password: password,
				}

				result = dbTrx.Create(&loginDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				authDatastore = datastoremodel.Auth{
					UserID: userDatastore.ID,
				}

				result = dbTrx.Create(&authDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				auth = authDatastore.ToDomain()

				currentPassword := password
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				body = fmt.Sprintf(`
				{
					"current_password":"%s",
					"new_password":"%s"
				}`,
					passwords.CurrentPassword, passwords.NewPassword)

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsInvalid",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = fake.Word()
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				auth = domainmodel.Auth{}

				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				body = fmt.Sprintf(`
					"current_password":"%s",
					"new_password":"%s"
				`,
					passwords.CurrentPassword, passwords.NewPassword)

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth
			},
			StatusCode: http.StatusBadRequest,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				authN = authpkg.New(ts.RSAKeys)

				auth = domainmodel.Auth{}

				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = securitypkg.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				body = fmt.Sprintf(`
				{
					"current_password":"%s",
					"new_password":"%s"
				}`,
					passwords.CurrentPassword, passwords.NewPassword)

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(dbTrx)

			userDatastoreRepository := userdatastorerepository.New(dbTrx)

			loginDatastoreRepository := logindatastorerepository.New(dbTrx)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "ChangePassword",
				Method:      http.MethodPost,
				Path:        "/change_password",
				HandlerFunc: authHandler.ChangePassword,
			}

			requestData := requesthttputilpkg.RequestData{
				Method:     route.Method,
				Target:     route.Path,
				Body:       body,
				ContextMap: contextMap,
			}

			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

			requesthttputilpkg.SetRequestContext(req, requestData.ContextMap)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedMessage := messagehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
				assert.NotEmpty(t, returnedMessage.Text)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}

func (ts *TestSuite) TestSignOut() {
	dbTrx := &gorm.DB{}

	var authN authpkg.IAuth

	userDatastore := datastoremodel.User{}

	loginDatastore := datastoremodel.Login{}

	authDatastore := datastoremodel.Auth{}

	auth := domainmodel.Auth{}

	contextMap := make(map[interface{}]interface{})

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningOut",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				userDatastore = datastoremodel.User{
					Username: username,
				}

				result := dbTrx.Create(&userDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				loginDatastore = datastoremodel.Login{
					UserID:   userDatastore.ID,
					Username: username,
					Password: password,
				}

				result = dbTrx.Create(&loginDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				authDatastore = datastoremodel.Auth{
					UserID: userDatastore.ID,
				}

				result = dbTrx.Create(&authDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				auth = authDatastore.ToDomain()

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsInvalid",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = fake.Word()
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

				authN = authpkg.New(ts.RSAKeys)

				auth = domainmodel.Auth{}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(dbTrx)

			userDatastoreRepository := userdatastorerepository.New(dbTrx)

			loginDatastoreRepository := logindatastorerepository.New(dbTrx)

			authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "SignOut",
				Method:      http.MethodPost,
				Path:        "/sign_out",
				HandlerFunc: authHandler.SignOut,
			}

			requestData := requesthttputilpkg.RequestData{
				Method:     route.Method,
				Target:     route.Path,
				ContextMap: contextMap,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			requesthttputilpkg.SetRequestContext(req, requestData.ContextMap)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedMessage := messagehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
				assert.NotEmpty(t, returnedMessage.Text)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}
