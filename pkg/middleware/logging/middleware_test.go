package logging_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	loggingmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestMiddlewareUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestLogging() {
	req := httptest.NewRequest(http.MethodGet, "/testing", nil)
	resprec := httptest.NewRecorder()

	handler := func(w http.ResponseWriter, r *http.Request) {}

	router := mux.NewRouter()
	router.Name("testing").
		Methods(http.MethodGet).
		Path("/testing").
		HandlerFunc(handler)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	loggedRouter.ServeHTTP(resprec, req)

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionAndLogsRequestsToOut",
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			loggingMiddleware := loggingmiddlewarepkg.Logging()

			testReq := httptest.NewRequest(http.MethodGet, "/testing", nil)
			testResprec := httptest.NewRecorder()
			logging := loggingMiddleware(handler)
			logging.ServeHTTP(testResprec, testReq)

			assert.Equal(t, resprec, testResprec)
			assert.Equal(t, req, testReq)
		})
	}
}
