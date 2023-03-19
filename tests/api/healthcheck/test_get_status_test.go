package healthcheck_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	healthcheckservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/healthcheck"
	healthcheckhandler "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/handler/healthcheck"
	requesthttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/request"
	responsehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/response"
	routehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/route"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestGetStatus() {
	db := &gorm.DB{}

	message := responsehttputilpkg.Message{}

	var connPool gorm.ConnPool

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInGettingTheStatus",
			SetUp: func(t *testing.T) {
				db = ts.DB

				message = responsehttputilpkg.Message{Text: "everything is up and running"}
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
				returnedMessage := responsehttputilpkg.Message{}
				err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
				assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v.", err))
				assert.Equal(t, message.Text, returnedMessage.Text)
			} else {
				assert.Equal(t, resprec.Code, tc.StatusCode)
			}

			tc.TearDown(t)
		})
	}
}
