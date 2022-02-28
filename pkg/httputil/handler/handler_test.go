package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	handlerhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHandlerUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestGetNotFoundHandler() {
	statusCode := 0

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInConfiguringNotFoundHandler",
			SetUp: func(t *testing.T) {
				statusCode = http.StatusNotFound
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			req := httptest.NewRequest(http.MethodGet, "/testing", nil)

			resprec := httptest.NewRecorder()

			handler := handlerhttputilpkg.GetNotFoundHandler()

			handler.ServeHTTP(resprec, req)

			assert.Equal(t, statusCode, resprec.Result().StatusCode)
		})
	}
}

func (ts *TestSuite) TestGetMethodNotAllowedHandler() {
	statusCode := 0

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInConfiguringMethodNotAllowedHandler",
			SetUp: func(t *testing.T) {
				statusCode = http.StatusMethodNotAllowed
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			req := httptest.NewRequest(http.MethodGet, "/testing", nil)

			resprec := httptest.NewRecorder()

			handler := handlerhttputilpkg.GetMethodNotAllowedHandler()

			handler.ServeHTTP(resprec, req)

			assert.Equal(t, statusCode, resprec.Result().StatusCode)
		})
	}
}
