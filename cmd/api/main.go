package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	_ "github.com/icaroribeiro/go-code-challenge-template/docs/api/swagger"
	authservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/auth"
	healthcheckservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/healthcheck"
	userservice "github.com/icaroribeiro/go-code-challenge-template/internal/application/service/user"
	authdatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/go-code-challenge-template/internal/infrastructure/datastore/repository/user"
	authhandler "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/handler/auth"
	healthcheckhandler "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/handler/healthcheck"
	userhandler "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/handler/user"
	authrouter "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/router/auth"
	healthcheckrouter "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/router/healthcheck"
	swaggerrouter "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/router/swagger"
	userrouter "github.com/icaroribeiro/go-code-challenge-template/internal/presentation/api/router/user"
	authpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/auth"
	datastorepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/datastore"
	envpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/env"
	adapterhttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/adapter"
	handlerhttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/handler"
	routehttputilpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/auth"
	dbtrxmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/dbtrx"
	loggingmiddlewarepkg "github.com/icaroribeiro/go-code-challenge-template/pkg/middleware/logging"
	securitypkg "github.com/icaroribeiro/go-code-challenge-template/pkg/security"
	serverpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/server"
	validatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator"
	passwordvalidatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator/password"
	usernamevalidatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator/username"
	uuidvalidatorpkg "github.com/icaroribeiro/go-code-challenge-template/pkg/validator/uuid"
	httpswaggerpkg "github.com/swaggo/http-swagger"
	validatorv2 "gopkg.in/validator.v2"
)

var (
	httpPort = envpkg.GetEnvWithDefaultValue("HTTP_PORT", "8080")

	publicKeyPath                  = envpkg.GetEnvWithDefaultValue("RSA_PUBLIC_KEY_PATH", "./configs/auth/rsa_keys/rsa.public")
	privateKeyPath                 = envpkg.GetEnvWithDefaultValue("RSA_PRIVATE_KEY_PATH", "./configs/auth/rsa_keys/rsa.private")
	tokenExpTimeInSecStr           = envpkg.GetEnvWithDefaultValue("TOKEN_EXP_TIME_IN_SEC", "120")
	timeBeforeTokenExpTimeInSecStr = envpkg.GetEnvWithDefaultValue("TIME_BEFORE_TOKEN_EXP_TIME_IN_SEC", "30")

	dbDriver   = envpkg.GetEnvWithDefaultValue("DB_DRIVER", "postgres")
	dbUser     = envpkg.GetEnvWithDefaultValue("DB_USER", "postgres")
	dbPassword = envpkg.GetEnvWithDefaultValue("DB_PASSWORD", "postgres")
	dbHost     = envpkg.GetEnvWithDefaultValue("DB_HOST", "localhost")
	dbPort     = envpkg.GetEnvWithDefaultValue("DB_PORT", "5433")
	dbName     = envpkg.GetEnvWithDefaultValue("DB_NAME", "db")
)

// @title New Go Code Challenge Template API
// @version 1.0
// @Description A REST API developed using Golang, Json Web Token and PostgreSQL database.
// @tag.name health check
// @tag.description It refers to the operation related to health check.
// @tag.name authentication
// @tag.description It refers to the operations related to authentication.
// @tag.name user
// @tag.description It refers to the operations related to user.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email icaroribeiro@hotmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @schemes http
// @SecurityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	httpPort := setupHttpPort()

	rsaKeys, err := setupRSAKeys()
	if err != nil {
		log.Panic(err.Error())
	}

	authN := authpkg.New(rsaKeys)

	tokenExpTimeInSec, err := strconv.Atoi(tokenExpTimeInSecStr)
	if err != nil {
		log.Panic(err.Error())
	}

	timeBeforeTokenExpTimeInSec, err := strconv.Atoi(timeBeforeTokenExpTimeInSecStr)
	if err != nil {
		log.Panic(err.Error())
	}

	dbConfig, err := setupDBConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	datastore, err := datastorepkg.New(dbConfig)
	if err != nil {
		log.Panic(err.Error())
	}
	defer datastore.Close()

	db := datastore.GetInstance()
	if db == nil {
		log.Panicf("The database instance is null")
	}

	if err = db.Error; err != nil {
		log.Panicf("Got error when acessing the database instance: %s", err.Error())
	}

	persistentAuthRepository := authdatastorerepository.New(db)
	persistentLoginRepository := logindatastorerepository.New(db)
	persistentUserRepository := userdatastorerepository.New(db)

	validationFuncs := map[string]validatorv2.ValidationFunc{
		"uuid":     uuidvalidatorpkg.Validate,
		"username": usernamevalidatorpkg.Validate,
		"password": passwordvalidatorpkg.Validate,
	}

	validator, err := validatorpkg.New(validationFuncs)
	if err != nil {
		log.Panic(err.Error())
	}

	security := securitypkg.New()

	healthCheckService := healthcheckservice.New(db)
	authService := authservice.New(persistentAuthRepository, persistentLoginRepository, persistentUserRepository,
		authN, security, validator, tokenExpTimeInSec)
	userService := userservice.New(persistentUserRepository, validator)

	swaggerHandler := httpswaggerpkg.WrapHandler
	healthCheckHandler := healthcheckhandler.New(healthCheckService)
	authHandler := authhandler.New(authService)
	userHandler := userhandler.New(userService)

	adapters := map[string]adapterhttputilpkg.Adapter{
		"loggingMiddleware":     loggingmiddlewarepkg.Logging(),
		"authMiddleware":        authmiddlewarepkg.Auth(db, authN),
		"authRenewalMiddleware": authmiddlewarepkg.AuthRenewal(db, authN, timeBeforeTokenExpTimeInSec),
		"dbTrxMiddleware":       dbtrxmiddlewarepkg.DBTrx(db),
	}

	routes := make(routehttputilpkg.Routes, 0)
	routes = append(routes, swaggerrouter.ConfigureRoutes(swaggerHandler, adapters)...)
	routes = append(routes, healthcheckrouter.ConfigureRoutes(healthCheckHandler, adapters)...)
	routes = append(routes, authrouter.ConfigureRoutes(authHandler, adapters)...)
	routes = append(routes, userrouter.ConfigureRoutes(userHandler, adapters)...)

	router := setupRouter(routes)

	server := serverpkg.New(fmt.Sprintf(":%s", httpPort), router)

	idleChan := make(chan struct{})

	go func() {
		waitForShutdown(*server)
		close(idleChan)
	}()

	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Panicf("%s", err.Error())
	}

	<-idleChan
}

// setupHttpPort is the function that configures the port address used by the server.
func setupHttpPort() string {
	return httpPort
}

// setupRSAKeys is the function that configures the RSA keys.
func setupRSAKeys() (authpkg.RSAKeys, error) {
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to read the RSA public key file: %s", err.Error())
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to parse the RSA public key: %s", err.Error())
	}

	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to read the RSA private key file: %s", err.Error())
	}

	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to parse the RSA private key: %s", err.Error())
	}

	return authpkg.RSAKeys{
		PublicKey:  rsaPublicKey,
		PrivateKey: rsaPrivateKey,
	}, nil
}

// setupDBConfig is the function that configures a map of parameters used to connect to the database.
func setupDBConfig() (map[string]string, error) {
	dbURL := ""

	dbConfig := map[string]string{
		"DRIVER":   dbDriver,
		"USER":     dbUser,
		"PASSWORD": dbPassword,
		"HOST":     dbHost,
		"PORT":     dbPort,
		"NAME":     dbName,
		"URL":      dbURL,
	}

	return dbConfig, nil
}

// setupRouter is the function that builds the router by arranging API routes.
func setupRouter(apiRoutes routehttputilpkg.Routes) *mux.Router {
	router := mux.NewRouter()

	methodNotAllowedHandler := handlerhttputilpkg.GetMethodNotAllowedHandler()
	router.MethodNotAllowedHandler = methodNotAllowedHandler

	notFoundHandler := handlerhttputilpkg.GetNotFoundHandler()
	router.NotFoundHandler = notFoundHandler

	for _, apiRoute := range apiRoutes {
		route := router.NewRoute()
		route.Name(apiRoute.Name)
		route.Methods(apiRoute.Method)

		if apiRoute.PathPrefix != "" {
			route.PathPrefix(apiRoute.PathPrefix)
		}

		if apiRoute.Path != "" {
			route.Path(apiRoute.Path)
		}

		route.HandlerFunc(apiRoute.HandlerFunc)
	}

	return router
}

// waitForShutdown is the function that waits for a signal to shutdown the server.
func waitForShutdown(server serverpkg.Server) {
	interruptChan := make(chan os.Signal, 1)

	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	<-interruptChan

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := server.Stop(ctx); err != nil && err != context.DeadlineExceeded {
		log.Panicf("%s", err.Error())
	}
}
