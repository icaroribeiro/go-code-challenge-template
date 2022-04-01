package healthcheck_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/application/service/healthcheck"
	healthcheckhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/http/presentation/handler/healthcheck"
	"github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestHealthCheckInteg(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestGetStatus() {
	db := &gorm.DB{}

	message := message.Message{}

	var connPool gorm.ConnPool

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingTheStatus",
			SetUp: func(t *testing.T) {
				db = ts.DB

				message = messagehttputilpkg.Message{Text: "everything is up and running"}
			},
			StatusCode: http.StatusOK,
			WantError:  false,
			TearDown:   func(t *testing.T) {},
		},
		{
			Context: "ItShouldFailIfTheDatabaseConnectionPoolIsInvalid",
			SetUp: func(t *testing.T) {
				connPool = ts.DB.ConnPool
				ts.DB.ConnPool = nil
				db = ts.DB
			},
			StatusCode: http.StatusInternalServerError,
			WantError:  true,
			TearDown: func(t *testing.T) {
				ts.DB.ConnPool = connPool
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			healthCheckService := healthcheckservice.New(db)
			healthCheckHandler := healthcheckhandler.New(healthCheckService)

			route := routehttputilpkg.Route{
				Name:        "GetStatus",
				Method:      "GET",
				Path:        "/status",
				HandlerFunc: healthCheckHandler.GetStatus,
			}

			requestData := requesthttputilpkg.RequestData{
				Method: route.Method,
				Target: route.Path,
			}

			reqBody := requesthttputilpkg.PrepareRequestBody(requestData.Body)

			req := httptest.NewRequest(requestData.Method, requestData.Target, reqBody)

			resprec := httptest.NewRecorder()

			router := mux.NewRouter()

			router.Name(route.Name).
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.HandlerFunc)

			router.ServeHTTP(resprec, req)

			if !tc.WantError {
				assert.Equal(t, resprec.Code, tc.StatusCode)
				returnedMessage := messagehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error %v.", err))
				assert.Equal(t, message.Text, returnedMessage.Text)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}
