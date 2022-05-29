package user_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/mockservice/user"
	userhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/handler/user"
	httppresentationmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/model"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	domainmodelfactory "github.com/icaroribeiro/new-go-code-challenge-template/tests/factory/core/domain/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestHandlerUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestGetAll() {
	user := domainmodel.User{}

	dbTrx := &gorm.DB{}
	dbTrx = nil

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				user = domainmodelfactory.NewUser(nil)

				returnArgs = ReturnArgs{
					{domainmodel.Users{user}, nil},
				}
			},
			StatusCode: http.StatusOK,
			WantError:  false,
		},
		{
			Context: "ItShouldFailIfAnErrorOccursWhenGettingAllUsers",
			SetUp: func(t *testing.T) {
				returnArgs = ReturnArgs{
					{domainmodel.Users{}, customerror.New("failed")},
				}
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			userService := new(usermockservice.Service)
			userService.On("WithDBTrx", dbTrx).Return(userService)
			userService.On("GetAll").Return(returnArgs[0]...)

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

			requesthttputilpkg.SetRequestHeaders(req, requestData.Headers)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedUsers := make(httppresentationmodel.Users, 0)
				err := json.NewDecoder(resprec.Body).Decode(&returnedUsers)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, user.ID, returnedUsers[0].ID)
				assert.Equal(t, user.Username, returnedUsers[0].Username)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}
		})
	}
}
