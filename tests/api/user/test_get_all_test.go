package user_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	userservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/user"
	datastoreentity "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/perentity"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/user"
	userhandler "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/handler/user"
	presentableentity "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/presentity"
	requesthttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/route"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestGetAll() {
	dbTrx := &gorm.DB{}

	persistentUser := datastoreentity.User{}

	user := presentableentity.User{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				username := fake.Username()

				persistentUser = datastoreentity.User{
					Username: username,
				}

				result := dbTrx.Create(&persistentUser)
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))

				domainUser := persistentUser.ToDomain()
				user.FromDomain(domainUser)
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown: func(t *testing.T) {
				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
		},
		{
			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
			SetUp: func(t *testing.T) {
				dbTrx = ts.DB.Begin()
				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error: %v.", dbTrx.Error))

				result := dbTrx.Rollback()
				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error: %v.", result.Error))
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown:   func(t *testing.T) {},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			persistentUserRepository := userdatastorerepository.New(dbTrx)
			userService := userservice.New(persistentUserRepository, ts.Validator)
			userHandler := userhandler.New(userService)

			route := routehttputilpkg.Route{
				Name:        "GetAllUsers",
				Method:      "GET",
				Path:        "/users",
				HandlerFunc: userHandler.GetAll,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
			}

			req := httptest.NewRequest(requestData.Method, requestData.Target, nil)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedUsers := presentableentity.Users{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedUsers)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, user.ID, returnedUsers[0].ID)
				assert.Equal(t, user.Username, returnedUsers[0].Username)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}
