package auth_test

// import (
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"net/http"
// 	"testing"

// 	fake "github.com/brianvoe/gofakeit/v5"
// 	authmodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/auth"
// 	authservicemock "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/servicemock/auth"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/security"
// 	authhandler "github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/handler/auth"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/httputil"
// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/gorm"
// )

// func TestHandler(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestSignUp() {
// 	credentials := security.Credentials{}

// 	body := ""

// 	dbTrx := &gorm.DB{}

// 	contextMap := make(map[interface{}]interface{})

// 	token := ""

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSigningUp",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"username":"%s",
// 					"password":"%s"
// 				}`,
// 					credentials.Username, credentials.Password)

// 				dbTrx = ts.DB.Begin()

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx

// 				token = fake.Word()

// 				returnArgs = ReturnArgs{
// 					{token, nil},
// 				}
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				dbTrx = nil

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx

// 				returnArgs = ReturnArgs{
// 					{"", nil},
// 				}

// 				body = ""
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 					"username":"%s",
// 					"password":"%s"
// 				`,
// 					credentials.Username, credentials.Password)

// 				dbTrx = ts.DB.Begin()

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx

// 				returnArgs = ReturnArgs{
// 					{"", nil},
// 				}
// 			},
// 			StatusCode: http.StatusBadRequest,
// 			WantError:  true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenRegisteringTheCredentials",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"username":"%s",
// 					"password":"%s"
// 				}`,
// 					credentials.Username, credentials.Password)

// 				dbTrx = ts.DB.Begin()

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx

// 				returnArgs = ReturnArgs{
// 					{"", errors.New("failed")},
// 				}

// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authService := new(authservicemock.ServiceMock)
// 			authService.On("WithDBTrx", dbTrx).Return(authService)
// 			authService.On("Register", credentials).Return(returnArgs[0]...)

// 			authHandler := authhandler.New(authService)

// 			route := httputil.Route{
// 				Name:        "SignUp",
// 				Method:      "POST",
// 				Path:        "/sign_up",
// 				HandlerFunc: authHandler.SignUp,
// 			}

// 			request := httputil.Request{
// 				Body:       body,
// 				ContextMap: contextMap,
// 			}

// 			resp := httputil.ExecuteRequest(route, request)

// 			if !tc.WantError {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 				returnedToken := httputil.Token{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedToken)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.NotEmpty(t, returnedToken.Text)
// 				assert.Equal(t, httputil.Token{Text: token}, returnedToken)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestSignIn() {
// 	credentials := security.Credentials{}

// 	body := ""

// 	dbTrx := &gorm.DB{}

// 	contextMap := make(map[interface{}]interface{})

// 	token := ""

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInLoggingIn",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"username":"%s",
// 					"password":"%s"
// 				}`,
// 					credentials.Username, credentials.Password)

// 				dbTrx = ts.DB.Begin()

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx

// 				token = fake.Word()

// 				returnArgs = ReturnArgs{
// 					{token, nil},
// 				}
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"username":"%s",
// 					"password":"%s"
// 				}`,
// 					credentials.Username, credentials.Password)

// 				dbTrx = nil

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx

// 				returnArgs = ReturnArgs{
// 					{"", nil},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 					"username":"%s",
// 					"password":"%s"
// 				`,
// 					credentials.Username, credentials.Password)

// 				dbTrx = ts.DB.Begin()

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx

// 				returnArgs = ReturnArgs{
// 					{"", nil},
// 				}
// 			},
// 			StatusCode: http.StatusBadRequest,
// 			WantError:  true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenLoggingIn",
// 			SetUp: func(t *testing.T) {
// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"username":"%s",
// 					"password":"%s"
// 				}`,
// 					credentials.Username, credentials.Password)

// 				dbTrx = ts.DB.Begin()

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx

// 				returnArgs = ReturnArgs{
// 					{"", errors.New("failed")},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authService := new(authservicemock.ServiceMock)
// 			authService.On("WithDBTrx", dbTrx).Return(authService)
// 			authService.On("LogIn", credentials).Return(returnArgs[0]...)

// 			authHandler := authhandler.New(authService)

// 			route := httputil.Route{
// 				Name:        "SignIn",
// 				Method:      "POST",
// 				Path:        "/sign_in",
// 				HandlerFunc: authHandler.SignIn,
// 			}

// 			request := httputil.Request{
// 				Body:       body,
// 				ContextMap: contextMap,
// 			}

// 			resp := httputil.ExecuteRequest(route, request)

// 			if !tc.WantError {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 				returnedToken := httputil.Token{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedToken)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.NotEmpty(t, returnedToken.Text)
// 				assert.Equal(t, httputil.Token{Text: token}, returnedToken)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestRefreshToken() {
// 	auth := authmodel.Auth{}

// 	contextMap := make(map[interface{}]interface{})

// 	token := ""

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInRefreshingTheToken",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth

// 				token = fake.Word()

// 				returnArgs = ReturnArgs{
// 					{token, nil},
// 				}
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 		},
// 		{
// 			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = nil

// 				token = fake.Word()

// 				returnArgs = ReturnArgs{
// 					{token, nil},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenRefreshingTheToken",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth

// 				returnArgs = ReturnArgs{
// 					{"", errors.New("failed")},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authService := new(authservicemock.ServiceMock)
// 			authService.On("RenewToken", auth).Return(returnArgs[0]...)

// 			authHandler := authhandler.New(authService)

// 			route := httputil.Route{
// 				Name:        "RefreshToken",
// 				Method:      "POST",
// 				Path:        "/refresh_token",
// 				HandlerFunc: authHandler.RefreshToken,
// 			}

// 			request := httputil.Request{
// 				ContextMap: contextMap,
// 			}

// 			resp := httputil.ExecuteRequest(route, request)

// 			if !tc.WantError {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 				returnedToken := httputil.Token{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedToken)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.NotEmpty(t, returnedToken.Text)
// 				assert.Equal(t, httputil.Token{Text: token}, returnedToken)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestChangePassword() {
// 	passwords := security.Passwords{}

// 	body := ""

// 	auth := authmodel.Auth{}

// 	contextMap := make(map[interface{}]interface{})

// 	message := ""

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInResettingThePassword",
// 			SetUp: func(t *testing.T) {
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"current_password":"%s",
// 					"new_password":"%s"
// 				}`,
// 					passwords.CurrentPassword, passwords.NewPassword)

// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth

// 				message = "the password has been updated successfully"

// 				returnArgs = ReturnArgs{
// 					{nil},
// 				}
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 		},
// 		{
// 			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
// 			SetUp: func(t *testing.T) {
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = nil

// 				returnArgs = ReturnArgs{
// 					{nil},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 		{
// 			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
// 			SetUp: func(t *testing.T) {
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				body = fmt.Sprintf(`
// 					"current_password":"%s",
// 					"new_password":"%s"
// 				`,
// 					passwords.CurrentPassword, passwords.NewPassword)

// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth

// 				returnArgs = ReturnArgs{
// 					{nil},
// 				}
// 			},
// 			StatusCode: http.StatusBadRequest,
// 			WantError:  true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenResettingThePassword",
// 			SetUp: func(t *testing.T) {
// 				currentPassword := fake.Password(true, true, true, false, false, 8)
// 				newPassword := fake.Password(true, true, true, false, false, 8)

// 				passwords = security.Passwords{
// 					CurrentPassword: currentPassword,
// 					NewPassword:     newPassword,
// 				}

// 				body = fmt.Sprintf(`
// 				{
// 					"current_password":"%s",
// 					"new_password":"%s"
// 				}`,
// 					passwords.CurrentPassword, passwords.NewPassword)

// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth

// 				returnArgs = ReturnArgs{
// 					{errors.New("failed")},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authService := new(authservicemock.ServiceMock)
// 			authService.On("ModifyPassword", auth.UserID.String(), passwords).Return(returnArgs[0]...)

// 			authHandler := authhandler.New(authService)

// 			route := httputil.Route{
// 				Name:        "ChangePassword",
// 				Method:      "POST",
// 				Path:        "/change_password",
// 				HandlerFunc: authHandler.ChangePassword,
// 			}

// 			request := httputil.Request{
// 				Body:       body,
// 				ContextMap: contextMap,
// 			}

// 			resp := httputil.ExecuteRequest(route, request)

// 			if !tc.WantError {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 				returnedMessage := httputil.Message{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedMessage)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.NotEmpty(t, returnedMessage.Text)
// 				assert.Equal(t, httputil.Message{Text: message}, returnedMessage)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}
// 		})
// 	}
// }

// func (ts *TestSuite) TestSignOut() {
// 	auth := authmodel.Auth{}

// 	contextMap := make(map[interface{}]interface{})

// 	message := ""

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSigningOut",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth

// 				message = "you have logged out successfully"

// 				returnArgs = ReturnArgs{
// 					{nil},
// 				}
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 		},
// 		{
// 			Context: "ItShouldFailIfItIsNotPossibleToGetTheAuthFromTheRequestContext",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = nil

// 				returnArgs = ReturnArgs{
// 					{nil},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenSigningOut",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				userID := uuid.NewV4()

// 				auth = authmodel.Auth{
// 					ID:     id,
// 					UserID: userID,
// 				}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth

// 				returnArgs = ReturnArgs{
// 					{errors.New("failed")},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authService := new(authservicemock.ServiceMock)
// 			authService.On("LogOut", auth.ID.String()).Return(returnArgs[0]...)

// 			authHandler := authhandler.New(authService)

// 			route := httputil.Route{
// 				Name:        "LogOut",
// 				Method:      "POST",
// 				Path:        "/sign_out",
// 				HandlerFunc: authHandler.SignOut,
// 			}

// 			request := httputil.Request{
// 				ContextMap: contextMap,
// 			}

// 			resp := httputil.ExecuteRequest(route, request)

// 			if !tc.WantError {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 				returnedMessage := httputil.Message{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedMessage)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.NotEmpty(t, returnedMessage.Text)
// 				assert.Equal(t, httputil.Message{Text: message}, returnedMessage)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}
// 		})
// 	}
// }
