package user_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/gorilla/mux"
	domainmodel "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/domain/model"
	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/mockservice/user"
	userhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/http/presentation/handler/user"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/customerror"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHandlerUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestGetAll() {
	user := domainmodel.User{}

	returnArgs := ReturnArgs{}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingAllUsers",
			SetUp: func(t *testing.T) {
				id := uuid.NewV4()
				username := fake.Username()

				user = domainmodel.User{
					ID:       id,
					Username: username,
				}

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

			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

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
				returnedUsers := make(domainmodel.Users, 0)
				err := json.NewDecoder(resprec.Body).Decode(&returnedUsers)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
				assert.Equal(t, returnedUsers[0].ID, user.ID)
				assert.Equal(t, returnedUsers[0].Username, user.Username)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}
		})
	}
}
