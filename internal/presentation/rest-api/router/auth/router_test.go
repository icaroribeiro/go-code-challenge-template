package auth_test

import (
	"net/http"
	"reflect"
	"runtime"
	"testing"

	authmockservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/core/ports/application/mockservice/auth"
	authhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/presentation/rest-api/handler/auth"
	authrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/presentation/rest-api/router/auth"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/auth"
	dbtrxmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/dbtrx"
	loggingmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

func TestRouterUnit(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) TestConfigureRoutes() {
	routes := routehttputilpkg.Routes{}

	db := &gorm.DB{}
	timeBeforeTokenExpTimeInSec := 10
	authN := authpkg.New(authpkg.RSAKeys{})

	authService := new(authmockservice.Service)
	authHandler := authhandler.New(authService)

	adapters := map[string]adapterhttputilpkg.Adapter{
		"loggingMiddleware":     loggingmiddlewarepkg.Logging(),
		"authMiddleware":        authmiddlewarepkg.Auth(db, authN),
		"authRenewalMiddleware": authmiddlewarepkg.AuthRenewal(db, authN, timeBeforeTokenExpTimeInSec),
		"dbTrxMiddleware":       dbtrxmiddlewarepkg.DBTrx(db),
	}

	ts.Cases = Cases{
		{
			Context: "ItShouldSucceedInConfiguringTheRoutes",
			SetUp: func(t *testing.T) {
				routes = routehttputilpkg.Routes{
					routehttputilpkg.Route{
						Name:   "SignUp",
						Method: http.MethodPost,
						Path:   "/sign_up",
						HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.SignUp).
							With(adapters["loggingMiddleware"], adapters["dbTrxMiddleware"]),
					},
					routehttputilpkg.Route{
						Name:   "SignIn",
						Method: http.MethodPost,
						Path:   "/sign_in",
						HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.SignIn).
							With(adapters["loggingMiddleware"], adapters["dbTrxMiddleware"]),
					},
					routehttputilpkg.Route{
						Name:   "RefreshToken",
						Method: http.MethodPost,
						Path:   "/refresh_token",
						HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.RefreshToken).
							With(adapters["loggingMiddleware"], adapters["authRenewalMiddleware"]),
					},
					routehttputilpkg.Route{
						Name:   "ChangePassword",
						Method: http.MethodPost,
						Path:   "/change_password",
						HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.ChangePassword).
							With(adapters["loggingMiddleware"], adapters["authMiddleware"]),
					},
					routehttputilpkg.Route{
						Name:   "SignOut",
						Method: http.MethodPost,
						Path:   "/sign_out",
						HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.SignOut).
							With(adapters["loggingMiddleware"], adapters["authMiddleware"]),
					},
				}
			},
		},
	}

	for _, tc := range ts.Cases {
		ts.T().Run(tc.Context, func(t *testing.T) {
			tc.SetUp(t)

			returnedRoutes := authrouter.ConfigureRoutes(authHandler, adapters)

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
