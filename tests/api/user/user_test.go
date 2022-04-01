package user_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"testing"

// 	fake "github.com/brianvoe/gofakeit/v5"
// 	userservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/user"
// 	usermodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/user"
// 	userdbmodel "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/model/user"
// 	userdsrepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/persistence/datastore/postgres/repository/user"
// 	userhandler "github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/handler/user"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/httputil"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"gorm.io/gorm"
// )

// func TestUser(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestGetAll() {
// 	dbTrx := &gorm.DB{}

// 	userdb := userdbmodel.User{}

// 	user := usermodel.User{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingAllUsers",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				username := fake.Username()

// 				userdb = userdbmodel.User{
// 					Username: username,
// 				}

// 				result := dbTrx.Create(&userdb)
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))

// 				user = userdb.ToDomain()
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 			TearDown: func(t *testing.T) {
// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
// 			},
// 		},
// 		{
// 			Context: "ItShouldFailIfTheDatabaseStateIsInconsistent",
// 			SetUp: func(t *testing.T) {
// 				dbTrx = ts.DB.Begin()
// 				assert.Nil(t, dbTrx.Error, fmt.Sprintf("Unexpected error %v.", dbTrx.Error))

// 				result := dbTrx.Rollback()
// 				assert.Nil(t, result.Error, fmt.Sprintf("Unexpected error %v.", result.Error))
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 			TearDown:   func(t *testing.T) {},
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			userDSRepository := userdsrepository.New(dbTrx)
// 			userService := userservice.New(userDSRepository, ts.Validator)
// 			userHandler := userhandler.New(userService)

// 			route := httputil.Route{
// 				Name:        "GetAllUsers",
// 				Method:      "GET",
// 				Path:        "/users",
// 				HandlerFunc: userHandler.GetAll,
// 			}

// 			request := httputil.Request{}

// 			resp := httputil.ExecuteRequest(route, request)

// 			if !tc.WantError {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 				returnedUsers := usermodel.Users{}
// 				err := json.NewDecoder(resp.Body).Decode(&returnedUsers)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, user.ID, returnedUsers[0].ID)
// 				assert.Equal(t, user.Username, returnedUsers[0].Username)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}

// 			tc.TearDown(t)
// 		})
// 	}
// }
