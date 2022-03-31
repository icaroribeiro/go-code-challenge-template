package user_test

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"testing"

// 	fake "github.com/brianvoe/gofakeit/v5"
// 	usermodel "github.com/icaroribeiro/go-code-challenge-template/internal/core/domain/model/user"
// 	userservicemock "github.com/icaroribeiro/go-code-challenge-template/internal/core/ports/application/servicemock/user"
// 	userhandler "github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/handler/user"
// 	"github.com/icaroribeiro/go-code-challenge-template/internal/interfaces/httputil"
// 	"github.com/icaroribeiro/go-code-challenge-template/pkg/customerror"
// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// func TestHandler(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }

// func (ts *TestSuite) TestGetAll() {
// 	user := usermodel.User{}

// 	returnArgs := ReturnArgs{}

// 	ts.Cases = Cases{
// 		{
// 			Context: "ItShouldSucceedInGettingAllUsers",
// 			SetUp: func(t *testing.T) {
// 				id := uuid.NewV4()
// 				username := fake.Username()

// 				user = usermodel.User{
// 					ID:       id,
// 					Username: username,
// 				}

// 				returnArgs = ReturnArgs{
// 					{usermodel.Users{user}, nil},
// 				}
// 			},
// 			StatusCode: http.StatusOK,
// 			WantError:  false,
// 		},
// 		{
// 			Context: "ItShouldFailIfAnErrorOccursWhenGettingAllUsers",
// 			SetUp: func(t *testing.T) {
// 				returnArgs = ReturnArgs{
// 					{usermodel.Users{}, customerror.New("failed")},
// 				}
// 			},
// 			StatusCode: http.StatusInternalServerError,
// 			WantError:  true,
// 		},
// 	}

// 	for _, tc := range ts.Cases {
// 		ts.T().Run(tc.Context, func(t *testing.T) {
// 			tc.SetUp(t)

// 			userService := new(userservicemock.ServiceMock)
// 			userService.On("GetAll").Return(returnArgs[0]...)

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
// 				returnedUsers := make(usermodel.Users, 0)
// 				err := json.NewDecoder(resp.Body).Decode(&returnedUsers)
// 				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
// 				assert.Equal(t, returnedUsers[0].ID, user.ID)
// 				assert.Equal(t, returnedUsers[0].Username, user.Username)
// 			} else {
// 				assert.Equal(t, resp.Code, tc.StatusCode)
// 			}
// 		})
// 	}
// }
