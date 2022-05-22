package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/mockservice/auth"
	authhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/handler/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	tokenhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/token"
	dbtrxmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/dbtrx"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/security"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestHandlerUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestSignUp() {
	credentials := security.Credentials{}

	body := ""

	driver := "postgres"
	db, mock := NewMockDB(driver)
	dbTrx := &gorm.DB{}

	//contextMap := make(map[interface{}]interface{})
	adapters := map[string]adapterhttputilpkg.Adapter{}

	token := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningUp",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = security.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				//var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				//contextMap[dbTrxCtxKey] = dbTrx

				dbTrx = db

				//mock.ExpectCommit()

				mock.ExpectBegin()
				mock.ExpectCommit()

				adapters["dbTrxMiddleware"] = dbtrxmiddlewarepkg.DBTrx(db)

				token = fake.Word()

				returnArgs = ReturnArgs{
					{token, nil},
				}
			},
			StatusCode: http.StatusOK,
			WantError:  false,
		},
		// {
		// 	Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
		// 	SetUp: func(t *testing.T) {
		// 		username := fake.Username()
		// 		password := fake.Password(true, true, true, false, false, 8)

		// 		credentials = security.Credentials{
		// 			Username: username,
		// 			Password: password,
		// 		}

		// 		dbTrx = nil

		// 		var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
		// 		contextMap[dbTrxKey] = dbTrx

		// 		returnArgs = ReturnArgs{
		// 			{"", nil},
		// 		}

		// 		body = ""
		// 	},
		// 	StatusCode: http.StatusInternalServerError,
		// 	WantError:  true,
		// },
		// {
		// 	Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
		// 	SetUp: func(t *testing.T) {
		// 		username := fake.Username()
		// 		password := fake.Password(true, true, true, false, false, 8)

		// 		credentials = security.Credentials{
		// 			Username: username,
		// 			Password: password,
		// 		}

		// 		body = fmt.Sprintf(`
		// 			"username":"%s",
		// 			"password":"%s"
		// 		`,
		// 			credentials.Username, credentials.Password)

		// 		dbTrx = db

		// 		var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
		// 		contextMap[dbTrxKey] = dbTrx

		// 		returnArgs = ReturnArgs{
		// 			{"", nil},
		// 		}
		// 	},
		// 	StatusCode: http.StatusBadRequest,
		// 	WantError:  true,
		// },
		// {
		// 	Context: "ItShouldFailIfAnErrorOccursWhenRegisteringTheCredentials",
		// 	SetUp: func(t *testing.T) {
		// 		username := fake.Username()
		// 		password := fake.Password(true, true, true, false, false, 8)

		// 		credentials = security.Credentials{
		// 			Username: username,
		// 			Password: password,
		// 		}

		// 		body = fmt.Sprintf(`
		// 		{
		// 			"username":"%s",
		// 			"password":"%s"
		// 		}`,
		// 			credentials.Username, credentials.Password)

		// 		dbTrx = db

		// 		var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
		// 		contextMap[dbTrxKey] = dbTrx

		// 		returnArgs = ReturnArgs{
		// 			{"", customerror.New("failed")},
		// 		}

		// 	},
		// 	StatusCode: http.StatusInternalServerError,
		// 	WantError:  true,
		// },
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authService := new(authmockservice.Service)
			authService.On("WithDBTrx", &dbTrx).Return(authService)
			authService.On("Register", credentials).Return(returnArgs[0]...)

			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:   "SignUp",
				Method: http.MethodPost,
				Path:   "/sign_up",
				HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.SignUp).
					With(adapters["dbTrxMiddleware"]),
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
				Body:   body,
				//ContextMap: contextMap,
			}

			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

			//requesthttputilpkg.SetRequestContext(req, requestData.ContextMap)

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
				assert.Equal(t, tokenhttputilpkg.Token{Text: token}, returnedToken)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}
		})
	}
}

func (ts *TestSuite) TestSignIn() {
	credentials := security.Credentials{}

	body := ""

	driver := "postgres"
	db, _ := NewMockDB(driver)
	dbTrx := &gorm.DB{}

	contextMap := make(map[interface{}]interface{})

	token := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInLoggingIn",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = security.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				dbTrx = db

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx

				token = fake.Word()

				returnArgs = ReturnArgs{
					{token, nil},
				}
			},
			StatusCode: http.StatusOK,
			WantError:  false,
		},
		{
			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = security.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				dbTrx = nil

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = security.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
					"username":"%s",
					"password":"%s"
				`,
					credentials.Username, credentials.Password)

				dbTrx = db

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx

				returnArgs = ReturnArgs{
					{"", nil},
				}
			},
			StatusCode: http.StatusBadRequest,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenLoggingIn",
			SetUp: func(t *testing.T) {
				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				credentials = security.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
				{
					"username":"%s",
					"password":"%s"
				}`,
					credentials.Username, credentials.Password)

				dbTrx = db

				var dbTrxKey requesthttputilpkg.ContextKeyType = "db_trx"
				contextMap[dbTrxKey] = dbTrx

				returnArgs = ReturnArgs{
					{"", customerror.New("failed")},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authService := new(authmockservice.Service)
			authService.On("WithDBTrx", dbTrx).Return(authService)
			authService.On("LogIn", credentials).Return(returnArgs[0]...)

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
				assert.Equal(t, tokenhttputilpkg.Token{Text: token}, returnedToken)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}
		})
	}
}

func (ts *TestSuite) TestRefreshToken() {
	auth := domainmodel.Auth{}

	dbTrx := &gorm.DB{}
	dbTrx = nil

	contextMap := make(map[interface{}]interface{})

	token := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRefreshingTheToken",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth

				token = fake.Word()

				returnArgs = ReturnArgs{
					{token, nil},
				}
			},
			StatusCode: http.StatusOK,
			WantError:  false,
		},
		{
			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = nil

				token = fake.Word()

				returnArgs = ReturnArgs{
					{token, nil},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenRefreshingTheToken",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth

				returnArgs = ReturnArgs{
					{"", customerror.New("failed")},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authService := new(authmockservice.Service)
			authService.On("WithDBTrx", dbTrx).Return(authService)
			authService.On("RenewToken", auth).Return(returnArgs[0]...)

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
				assert.Equal(t, tokenhttputilpkg.Token{Text: token}, returnedToken)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}
		})
	}
}

func (ts *TestSuite) TestChangePassword() {
	passwords := security.Passwords{}

	body := ""

	auth := domainmodel.Auth{}

	dbTrx := &gorm.DB{}
	dbTrx = nil

	contextMap := make(map[interface{}]interface{})

	message := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInResettingThePassword",
			SetUp: func(t *testing.T) {
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				body = fmt.Sprintf(`
				{
					"current_password":"%s",
					"new_password":"%s"
				}`,
					passwords.CurrentPassword, passwords.NewPassword)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth

				message = "the password has been updated successfully"

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			StatusCode: http.StatusOK,
			WantError:  false,
		},
		{
			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
			SetUp: func(t *testing.T) {
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = nil

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
			SetUp: func(t *testing.T) {
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				body = fmt.Sprintf(`
					"current_password":"%s",
					"new_password":"%s"
				`,
					passwords.CurrentPassword, passwords.NewPassword)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			StatusCode: http.StatusBadRequest,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenResettingThePassword",
			SetUp: func(t *testing.T) {
				currentPassword := fake.Password(true, true, true, false, false, 8)
				newPassword := fake.Password(true, true, true, false, false, 8)

				passwords = security.Passwords{
					CurrentPassword: currentPassword,
					NewPassword:     newPassword,
				}

				body = fmt.Sprintf(`
				{
					"current_password":"%s",
					"new_password":"%s"
				}`,
					passwords.CurrentPassword, passwords.NewPassword)

				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authService := new(authmockservice.Service)
			authService.On("WithDBTrx", dbTrx).Return(authService)
			authService.On("ModifyPassword", auth.UserID.String(), passwords).Return(returnArgs[0]...)

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
				assert.Equal(t, messagehttputilpkg.Message{Text: message}, returnedMessage)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}
		})
	}
}

func (ts *TestSuite) TestSignOut() {
	auth := domainmodel.Auth{}

	dbTrx := &gorm.DB{}
	dbTrx = nil

	contextMap := make(map[interface{}]interface{})

	message := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningOut",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth

				message = "you have logged out successfully"

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			StatusCode: http.StatusOK,
			WantError:  false,
		},
		{
			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = nil

				returnArgs = ReturnArgs{
					{nil},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenSigningOut",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				userID := uuid.NewV4()

				auth = domainmodel.Auth{
					ID:     id,
					UserID: userID,
				}

				var authDetailsKey requesthttputilpkg.ContextKeyType = "auth_details"
				contextMap[authDetailsKey] = auth

				returnArgs = ReturnArgs{
					{customerror.New("failed")},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authService := new(authmockservice.Service)
			authService.On("WithDBTrx", dbTrx).Return(authService)
			authService.On("LogOut", auth.ID.String()).Return(returnArgs[0]...)

			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "LogOut",
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
				assert.Equal(t, messagehttputilpkg.Message{Text: message}, returnedMessage)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}
		})
	}
}
