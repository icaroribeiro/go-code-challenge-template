package swagger

import (
	"net/http"

	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	loggingmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/logging"
	httpswaggerpkg "github.com/swaggo/http-swagger"
)

// ConfigureRoutes is the function that arranges the swagger's routes.
func ConfigureRoutes() routehttputilpkg.Routes {
	swaggerHandler := httpswaggerpkg.WrapHandler

	loggingMiddleware := loggingmiddlewarepkg.Logging()

	return routehttputilpkg.Routes{
		routehttputilpkg.Route{
			Name:       "Swagger",
			Method:     http.MethodGet,
			PathPrefix: "/swagger",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(swaggerHandler).
				With(loggingMiddleware),
		},
	}
}
