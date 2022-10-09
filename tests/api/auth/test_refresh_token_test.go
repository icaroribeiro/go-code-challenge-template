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
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/entity"
	datastoremodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/entity"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/user"
	authhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/presentation/rest-api/handler/auth"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	tokenhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/token"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/auth"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestRefreshToken() {
	dbTrx := &gorm.DB{}

	var authN authpkg.IAuth

	authDatastore := datastoremodel.Auth{}
	loginDatastore := datastoremodel.Login{}
	userDatastore := datastoremodel.User{}

	authDetailsCtxValue := domainmodel.Auth{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInRefreshingTheToken",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				username := fake.Username()
				password := fake.Password(true, true, true, false, false, 8)

				userDatastore = datastoremodel.User{
					Username: username,
				}

				result := dbTrx.Create(&userDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				loginDatastore = datastoremodel.Login{
					UserID:   userDatastore.ID,
					Username: username,
					Password: password,
				}

				result = dbTrx.Create(&loginDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				authDatastore = datastoremodel.Auth{
					UserID: userDatastore.ID,
				}

				result = dbTrx.Create(&authDatastore)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				authDetailsCtxValue = authDatastore.ToDomain()
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheAuthDetailsFromTheRequestContextIsEmpty",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				authN = authpkg.New(ts.RSAKeys)

				authDetailsCtxValue = domainmodel.Auth{}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			authDatastoreRepository := authdatastorerepository.New(dbTrx)
			loginDatastoreRepository := logindatastorerepository.New(dbTrx)
			userDatastoreRepository := userdatastorerepository.New(dbTrx)

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
				Method: route.Method,
				Target: route.Path,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			ctx := req.Context()
			ctx = authmiddlewarepkg.NewContext(ctx, authDetailsCtxValue)
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
