package healthcheck

import (
	"net/http"

	healthcheckhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/presentation/api/handler/healthcheck"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
)

// ConfigureRoutes is the function that arranges the healthcheck's routes.
func ConfigureRoutes(healthCheckHandler healthcheckhandler.IHandler, adapters map[string]adapterhttputilpkg.Adapter) routehttputilpkg.Routes {
	return routehttputilpkg.Routes{
		routehttputilpkg.Route{
			Name:   "GetStatus",
			Method: http.MethodGet,
			Path:   "/status",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(healthCheckHandler.GetStatus).
				With(adapters["loggingMiddleware"]),
		},
	}
}
