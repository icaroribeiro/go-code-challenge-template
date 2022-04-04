package auth

import (
	authhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/handler/auth"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
)

// ConfigureRoutes is the function that arranges the auth's routes.
func ConfigureRoutes(authHandler authhandler.IHandler, adapters map[string]adapterhttputilpkg.Adapter) routehttputilpkg.Routes {
	return routehttputilpkg.Routes{
		routehttputilpkg.Route{
			Name:   "SignUp",
			Method: "POST",
			Path:   "/sign_up",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.SignUp).
				With(adapters["loggingMiddleware"], adapters["dbTrxMiddleware"]),
		},
		routehttputilpkg.Route{
			Name:   "SignIn",
			Method: "POST",
			Path:   "/sign_in",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.SignIn).
				With(adapters["loggingMiddleware"], adapters["dbTrxMiddleware"]),
		},
		routehttputilpkg.Route{
			Name:   "RefreshToken",
			Method: "POST",
			Path:   "/refresh_token",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.RefreshToken).
				With(adapters["loggingMiddleware"], adapters["authRenewalMiddleware"]),
		},
		routehttputilpkg.Route{
			Name:   "ChangePassword",
			Method: "POST",
			Path:   "/change_password",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.ChangePassword).
				With(adapters["loggingMiddleware"], adapters["authMiddleware"]),
		},
		routehttputilpkg.Route{
			Name:   "SignOut",
			Method: "POST",
			Path:   "/sign_out",
			HandlerFunc: adapterhttputilpkg.AdaptFunc(authHandler.SignOut).
				With(adapters["loggingMiddleware"], adapters["authMiddleware"]),
		},
	}
}
