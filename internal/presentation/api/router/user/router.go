package user

import (
	userhandler "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/handler/user"
	adapterhttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/route"
)

// ConfigureRoutes is the function that arranges the user's routes.
func ConfigureRoutes(userHandler userhandler.IHandler, adapters map[string]adapterhttputilpkg.Adapter) routehttputilpkg.Routes {
	return routehttputilpkg.Routes{
		routehttputilpkg.Route{
			Name:   "GetAllUsers",
			Method: "GET",
			Path:   "/users",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(userHandler.GetAll).
				With(adapters["loggingMiddleware"], adapters["authMiddleware"]),
		},
	}
}
