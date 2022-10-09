package auth_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/mockservice/auth"
	authhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/presentation/rest-api/handler/auth"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	tokenhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/token"
	dbtrxmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/dbtrx"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/security"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestSignIn() {
	credentials := securitypkg.Credentials{}

	body := ""

	driver := "postgres"
	db, _ := NewMockDB(driver)

	dbTrxCtxValue := &gorm.DB{}

	tokenString := ""

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInSigningIn",
			SetUp: func(t *testing.T) {
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

				dbTrxCtxValue = db

				tokenString = fake.Word()

				returnArgs = ReturnArgs{
					{tokenString, nil},
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

				dbTrxCtxValue = nil

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

				credentials = securitypkg.Credentials{
					Username: username,
					Password: password,
				}

				body = fmt.Sprintf(`
					"username":"%s",
					"password":"%s"
				`,
					credentials.Username, credentials.Password)

				dbTrxCtxValue = db

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

				dbTrxCtxValue = db

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
			authService.On("WithDBTrx", dbTrxCtxValue).Return(authService)
			authService.On("LogIn", credentials).Return(returnArgs[0]...)

			authHandler := authhandler.New(authService)

			route := routehttputilpkg.Route{
				Name:        "SignIn",
				Method:      http.MethodPost,
				Path:        "/sign_in",
				HandlerFunc: authHandler.SignIn,
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
				assert.Equal(t, tokenhttputilpkg.Token{Text: tokenString}, returnedToken)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}
		})
	}
}
