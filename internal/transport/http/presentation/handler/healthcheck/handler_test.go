package healthcheck_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	healthcheckhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/http/presentation/handler/healthcheck"
	messagehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/message"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHandlerUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestGetStatus() {
	text := "everything is up and running"

	message := messagehttputilpkg.Message{Text: text}

	ts.Cases = Cases{
		{
			Context:    "ItShouldSucceedInGettingStatus",
			SetUp:      func(t *testing.T) {},
			StatusCode: http.StatusOK,
			WantError:  false,
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			healthCheckHandler := healthcheckhandler.New()

			route := routehttputilpkg.Route{
				Name:        "GetStatus",
				Method:      http.MethodGet,
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

			assert.Equal(t, resprec.Code, tc.StatusCode)
			returnedMessage := messagehttputilpkg.Message{}
			err := json.NewDecoder(resprec.Body).Decode(&returnedMessage)
			assert.Nil(t, err, fmt.Sprintf("Unexpected error: %v", err))
			assert.Equal(t, returnedMessage.Text, message.Text)
		})
	}
}
