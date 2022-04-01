package auth_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"go/token"
// 	"net/http"
// 	"testing"

// 	fake "github.com/brianvoe/gofakeit/v5"
// 	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/auth"
// 	authmodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/auth"
// 	loginmodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/login"
// 	usermodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/user"
// 	authinfra "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/auth"
// 	authdbmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/model/auth"
// 	logindbmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/model/login"
// 	userdbmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/model/user"
// 	authdsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/postgres/repository/auth"
// 	logindsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/postgres/repository/login"
// 	userdsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/postgres/repository/user"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/security"
// 	authhandler "github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/handler/auth"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/httputil"
// 	"github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/message"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/gorm"
// )

// func TestAuth(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestSignUp() {
// 	dbTrx := &gorm.DB{}

// 	var authInfra authinfra.IAuth

// 	credentials := security.Credentials{}

// 	body := ""

// 	contextMap := make(map[interface{}]interface{})

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSigningUp",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

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

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
// 			SetUp: func(t *testing.T) {
// 				authInfra = authinfra.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = ""

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = nil
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

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

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx
// 			},
// 			StatusCode: http.StatusBadRequest,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

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

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDSRepository := authdsrepository.New(dbTrx)

// 			userDSRepository := userdsrepository.New(dbTrx)

// 			loginDSRepository := logindsrepository.New(dbTrx)

// 			authService := authservice.New(authDSRepository, userDSRepository, loginDSRepository, authInfra,
// 				ts.TokenExpTimeInSec, ts.Security, ts.Validator)

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
// 				returnedToken := token.Token{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedToken)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.NotEmpty(t, returnedToken.Text)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestSignIn() {
// 	dbTrx := &gorm.DB{}

// 	var authInfra authinfra.IAuth

// 	userdb := userdbmodel.User{}

// 	logindb := logindbmodel.Login{}

// 	credentials := security.Credentials{}

// 	body := ""

// 	contextMap := make(map[interface{}]interface{})

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSigningIn",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				userdb = userdbmodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userdb)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				logindb = logindbmodel.Login{
// 					UserID:   userdb.ID,
// 					Username: username,
// 					Password: password,
// 				}

// 				result = dbTrx.Create(&logindb)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

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

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseTransactionFromTheRequestContextIsNull",
// 			SetUp: func(t *testing.T) {
// 				authInfra = authinfra.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				credentials = security.Credentials{
// 					Username: username,
// 					Password: password,
// 				}

// 				body = ""

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = nil
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

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

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx
// 			},
// 			StatusCode: http.StatusBadRequest,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

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

// 				var dbTrxKey httputil.ContextKeyType = "db_trx"
// 				contextMap[dbTrxKey] = dbTrx
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDSRepository := authdsrepository.New(dbTrx)

// 			userDSRepository := userdsrepository.New(dbTrx)

// 			loginDSRepository := logindsrepository.New(dbTrx)

// 			authService := authservice.New(authDSRepository, userDSRepository, loginDSRepository,
// 				authInfra, ts.TokenExpTimeInSec, ts.Security, ts.Validator)

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
// 				returnedToken := token.Token{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedToken)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.NotEmpty(t, returnedToken.Text)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestRefreshToken() {
// 	dbTrx := &gorm.DB{}

// 	var authInfra authinfra.IAuth

// 	user := usermodel.User{}

// 	login := loginmodel.Login{}

// 	auth := authmodel.Auth{}

// 	contextMap := make(map[interface{}]interface{})

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInRefreshingTheToken",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				user = usermodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&user)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				login = loginmodel.Login{
// 					UserID:   user.ID,
// 					Username: username,
// 					Password: password,
// 				}

// 				result = dbTrx.Create(&login)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				auth = authmodel.Auth{
// 					UserID: user.ID,
// 				}

// 				result = dbTrx.Create(&auth)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsInvalid",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = fake.Word()
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDSRepository := authdsrepository.New(dbTrx)

// 			userDSRepository := userdsrepository.New(dbTrx)

// 			loginDSRepository := logindsrepository.New(dbTrx)

// 			authService := authservice.New(authDSRepository, userDSRepository, loginDSRepository,
// 				authInfra, ts.TokenExpTimeInSec, ts.Security, ts.Validator)

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
// 				returnedToken := token.Token{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedToken)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.NotEmpty(t, returnedToken.Text)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestChangePassword() {
// 	dbTrx := &gorm.DB{}

// 	var authInfra authinfra.IAuth

// 	userdb := userdbmodel.User{}

// 	logindb := logindbmodel.Login{}

// 	authdb := authdbmodel.Auth{}

// 	auth := authmodel.Auth{}

// 	passwords := security.Passwords{}

// 	body := ""

// 	contextMap := make(map[interface{}]interface{})

// 	message := message.Message{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInResettingThePassword",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				userdb = userdbmodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userdb)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				logindb = logindbmodel.Login{
// 					UserID:   userdb.ID,
// 					Username: username,
// 					Password: password,
// 				}

// 				result = dbTrx.Create(&logindb)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				authdb = authdbmodel.Auth{
// 					UserID: userdb.ID,
// 				}

// 				result = dbTrx.Create(&authdb)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				auth = authdb.ToDomain()

// 				currentPassword := password
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

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth

// 				message = message.Message{Text: "the password has been updated successfully"}
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsInvalid",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = fake.Word()
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				auth = authmodel.Auth{}

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

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth
// 			},
// 			StatusCode: http.StatusBadRequest,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				auth = authmodel.Auth{}

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

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDSRepository := authdsrepository.New(dbTrx)

// 			userDSRepository := userdsrepository.New(dbTrx)

// 			loginDSRepository := logindsrepository.New(dbTrx)

// 			authService := authservice.New(authDSRepository, userDSRepository, loginDSRepository,
// 				authInfra, ts.TokenExpTimeInSec, ts.Security, ts.Validator)

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
// 				returnedMessage := message.Message{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedMessage)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, message.Text, returnedMessage.Text)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }

// func (ts *TestSuite) TestSignOut() {
// 	dbTrx := &gorm.DB{}

// 	var authInfra authinfra.IAuth

// 	userdb := userdbmodel.User{}

// 	logindb := logindbmodel.Login{}

// 	authdb := authdbmodel.Auth{}

// 	auth := authmodel.Auth{}

// 	contextMap := make(map[interface{}]interface{})

// 	message := message.Message{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInSigningOut",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				username := fake.Username()
// 				password := fake.Password(true, true, true, false, false, 8)

// 				userdb = userdbmodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userdb)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				logindb = logindbmodel.Login{
// 					UserID:   userdb.ID,
// 					Username: username,
// 					Password: password,
// 				}

// 				result = dbTrx.Create(&logindb)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				authdb = authdbmodel.Auth{
// 					UserID: userdb.ID,
// 				}

// 				result = dbTrx.Create(&authdb)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				auth = authdb.ToDomain()

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth

// 				message = message.Message{Text: "you have logged out successfully"}
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsInvalid",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = fake.Word()
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				authInfra = authinfra.New(ts.RSAKeys)

// 				auth = authmodel.Auth{}

// 				var authDetailsKey httputil.ContextKeyType = "auth_details"
// 				contextMap[authDetailsKey] = auth
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			authDSRepository := authdsrepository.New(dbTrx)

// 			userDSRepository := userdsrepository.New(dbTrx)

// 			loginDSRepository := logindsrepository.New(dbTrx)

// 			authService := authservice.New(authDSRepository, userDSRepository, loginDSRepository,
// 				authInfra, ts.TokenExpTimeInSec, ts.Security, ts.Validator)

// 			authHandler := authhandler.New(authService)

// 			route := httputil.Route{
// 				Name:        "SignOut",
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
// 				returnedMessage := message.Message{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedMessage)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, message.Text, returnedMessage.Text)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }
