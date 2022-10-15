package logging_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	requesthttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/request"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	loggingmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/logging"
	"github.com/stretchr/testify/assert"
)

func (ts *TestSuite) TestLogging() {
	handler := func(w http.ResponseWriter, r *http.Request) {}

	route := routehttputilpkg.Route{
		Name:        "Testing",
		Method:      http.MethodGet,
		Path:        "/testing",
		HandlerFunc: handler,
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

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	loggedRouter.ServeHTTP(resprec, req)

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInWrappingAFunctionWithLoggingMiddleware",
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			loggingMiddleware := loggingmiddlewarepkg.Logging()

			testReq := httptest.NewRequest(requestData.Method, requestData.Target, nil)
			testResprec := httptest.NewRecorder()
			logging := loggingMiddleware(handler)
			logging.ServeHTTP(testResprec, testReq)

			assert.Equal(t, resprec, testResprec)
			assert.Equal(t, req, testReq)
		})
	}
}
