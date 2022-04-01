package healthcheck_test

import (
	"net/http"
	"reflect"
	"runtime"
	"testing"

	healthcheckmockservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/mockservice/healthcheck"
	healthcheckhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/http/presentation/handler/healthcheck"
	healthcheckrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/http/router/healthcheck"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	loggingmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestRouterUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestConfigureRoutes() {
	routes := routehttputilpkg.Routes{}

	healthCheckService := new(healthcheckmockservice.Service)
	healthCheckHandler := healthcheckhandler.New(healthCheckService)

	adapters := map[string]adapterhttputilpkg.Adapter{
		"loggingMiddleware": loggingmiddlewarepkg.Logging(),
	}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInConfiguringTheRoutes",
			SetUp: func(t *testing.T) {
				routes = routehttputilpkg.Routes{
					routehttputilpkg.Route{
						Name:   "GetStatus",
						Method: http.MethodGet,
						Path:   "/status",
						HandlerFunc: adapterhttputilpkg.AdaptFunc(healthCheckHandler.GetStatus).
							With(adapters["loggingMiddleware"]),
					},
				}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedRoutes := healthcheckrouter.ConfigureRoutes(healthCheckHandler, adapters)

			assert.Equal(t, len(routes), len(returnedRoutes))

			for i := range routes {
				assert.Equal(t, routes[i].Name, returnedRoutes[i].Name)
				assert.Equal(t, routes[i].Method, returnedRoutes[i].Method)
				assert.Equal(t, routes[i].Path, returnedRoutes[i].Path)
				handlerFunc1 := runtime.FuncForPC(reflect.ValueOf(routes[i].HandlerFunc).Pointer()).Name()
				handlerFunc2 := runtime.FuncForPC(reflect.ValueOf(returnedRoutes[i].HandlerFunc).Pointer()).Name()
				assert.Equal(t, handlerFunc1, handlerFunc2)
			}
		})
	}
}
