package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/auth"
	authdatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/user"
	authhandler "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/handler/auth"
	authpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/auth"
	requesthttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/route"
	tokenhttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/token"
	dbtrxmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/dbtrx"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template/pkg/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestSignUp() {
	dbTrx := &gorm.DB{}

	var authN authpkg.IAuth

	credentials := securitypkg.Credentials{}

	body := ""

	dbTrxCtxValue := &gorm.DB{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningUp",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

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

				dbTrxCtxValue = dbTrx
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
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

				dbTrxCtxValue = nil
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheRequestBodyIsAnImproperlyFormattedJsonString",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

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

				dbTrxCtxValue = dbTrx
			},
			StatusCode: http.StatusBadRequest,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

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

				dbTrxCtxValue = dbTrx
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			persistentAuthRepository := authdatastorerepository.New(dbTrx)
			persistentLoginRepository := logindatastorerepository.New(dbTrx)
			persistentUserRepository := userdatastorerepository.New(dbTrx)

			authService := authservice.New(persistentAuthRepository, persistentLoginRepository, persistentUserRepository,
				authN, ts.Security, ts.Validator, ts.TokenExpTimeInSec)
			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "SignUp",
				Method:      http.MethodPost,
				Path:        "/sign_up",
				HandlerFunc: authHandler.SignUp,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
				Body:   body,
			}

			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

			ctx := req.Context()
			ctx = dbtrxmiddlewarepkg.NewContext(ctx, dbTrxCtxValue)
			req = req.WithContext(ctx)

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
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.NotEmpty(t, returnedToken.Text)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}
