package user_test

import (
	"reflect"
	"runtime"
	"testing"

	usermockservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/mockservice/user"
	userhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/presentation/rest-api/handler/user"
	userrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/presentation/rest-api/router/user"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/auth"
	loggingmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/logging"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func (ts *TestSuite) TestConfigureRoutes() {
	routes := routehttputilpkg.Routes{}

	db := &gorm.DB{}
	authN := authpkg.New(authpkg.RSAKeys{})

	userService := new(usermockservice.Service)
	userHandler := userhandler.New(userService)

	adapters := map[string]adapterhttputilpkg.Adapter{
		"loggingMiddleware": loggingmiddlewarepkg.Logging(),
		"authMiddleware":    authmiddlewarepkg.Auth(db, authN),
	}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInConfiguringTheRoutes",
			SetUp: func(t *testing.T) {
				routes = routehttputilpkg.Routes{
					routehttputilpkg.Route{
						Name:   "GetAllUsers",
						Method: "GET",
						Path:   "/users",
						HandlerFunc: adapterhttputilpkg.AdaptFunc(userHandler.GetAll).
							With(adapters["loggingMiddleware"], adapters["authMiddleware"]),
					},
				}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedRoutes := userrouter.ConfigureRoutes(userHandler, adapters)

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
