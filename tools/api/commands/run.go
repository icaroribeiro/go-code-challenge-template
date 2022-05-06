package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/icaroribeiro/new-go-code-challenge-template/docs/api/swagger"
	authservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/application/service/auth"
	healthcheckservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/application/service/healthcheck"
	userservice "github.com/icaroribeiro/new-go-code-challenge-template/internal/application/service/user"
	authdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/auth"
	logindatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/login"
	userdatastorerepository "github.com/icaroribeiro/new-go-code-challenge-template/internal/infrastructure/storage/datastore/repository/user"
	authhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/handler/auth"
	healthcheckhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/handler/healthcheck"
	userhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/presentation/handler/user"
	authrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/router/auth"
	healthcheckrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/router/healthcheck"
	swaggerrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/router/swagger"
	userrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/router/user"
	authpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/auth"
	datastorepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/datastore"
	envpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/env"
	adapterhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/adapter"
	handlerhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/handler"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	authmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/auth"
	dbtrxmiddleware "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/dbtrx"
	loggingmiddlewarepkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/middleware/logging"
	securitypkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/security"
	serverpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/server"
	validatorpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator"
	passwordvalidator "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator/password"
	usernamevalidator "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator/username"
	uuidvalidator "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator/uuid"
	"github.com/spf13/cobra"
	httpswaggerpkg "github.com/swaggo/http-swagger"
	validatorv2 "gopkg.in/validator.v2"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the API",
	Run:   execRunCmd,
}

var (
	deploy = envpkg.GetEnvWithDefaultValue("DEPLOY", "NO")

	httpPort = envpkg.GetEnvWithDefaultValue("HTTP_PORT", "8080")

	publicKeyPath                  = envpkg.GetEnvWithDefaultValue("RSA_PUBLIC_KEY_PATH", "./configs/auth/rsa_keys/rsa.public")
	privateKeyPath                 = envpkg.GetEnvWithDefaultValue("RSA_PRIVATE_KEY_PATH", "./configs/auth/rsa_keys/rsa.private")
	tokenExpTimeInSecStr           = envpkg.GetEnvWithDefaultValue("TOKEN_EXP_TIME_IN_SEC", "600")
	timeBeforeTokenExpTimeInSecStr = envpkg.GetEnvWithDefaultValue("TIME_BEFORE_TOKEN_EXP_TIME_IN_SEC", "60")

	dbDriver   = envpkg.GetEnvWithDefaultValue("DB_DRIVER", "postgres")
	dbUser     = envpkg.GetEnvWithDefaultValue("DB_USER", "postgres")
	dbPassword = envpkg.GetEnvWithDefaultValue("DB_PASSWORD", "postgres")
	dbHost     = envpkg.GetEnvWithDefaultValue("DB_HOST", "localhost")
	dbPort     = envpkg.GetEnvWithDefaultValue("DB_PORT", "5432")
	dbName     = envpkg.GetEnvWithDefaultValue("DB_NAME", "db")
)

func execRunCmd(cmd *cobra.Command, args []string) {
	tcpAddress := setupTcpAddress()

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

	authDatastoreRepository := authdatastorerepository.New(db)
	loginDatastoreRepository := logindatastorerepository.New(db)
	userDatastoreRepository := userdatastorerepository.New(db)

	validationFuncs := map[string]validatorv2.ValidationFunc{
		"uuid":     uuidvalidator.Validate,
		"username": usernamevalidator.Validate,
		"password": passwordvalidator.Validate,
	}

	validator, err := validatorpkg.New(validationFuncs)
	if err != nil {
		log.Panic(err.Error())
	}

	security := securitypkg.New()

	healthCheckService := healthcheckservice.New(db)
	authService := authservice.New(authDatastoreRepository, loginDatastoreRepository, userDatastoreRepository,
		authN, security, validator, tokenExpTimeInSec)
	userService := userservice.New(userDatastoreRepository, validator)

	swaggerHandler := httpswaggerpkg.WrapHandler
	healthCheckHandler := healthcheckhandler.New(healthCheckService)
	authHandler := authhandler.New(authService)
	userHandler := userhandler.New(userService)

	adapters := map[string]adapterhttputilpkg.Adapter{
		"loggingMiddleware":     loggingmiddlewarepkg.Logging(),
		"authMiddleware":        authmiddlewarepkg.Auth(db, authN),
		"authRenewalMiddleware": authmiddlewarepkg.AuthRenewal(db, authN, timeBeforeTokenExpTimeInSec),
		"dbTrxMiddleware":       dbtrxmiddleware.DBTrx(db),
	}

	routes := make(routehttputilpkg.Routes, 0)
	routes = append(routes, swaggerrouter.ConfigureRoutes(swaggerHandler, adapters)...)
	routes = append(routes, healthcheckrouter.ConfigureRoutes(healthCheckHandler, adapters)...)
	routes = append(routes, authrouter.ConfigureRoutes(authHandler, adapters)...)
	routes = append(routes, userrouter.ConfigureRoutes(userHandler, adapters)...)

	router := setupRouter(routes)

	server := serverpkg.New(tcpAddress, router)

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

// setupTcpAddress is the function that configures the tcp address used by the server.
func setupTcpAddress() string {
	if deploy == "YES" {
		if httpPort = os.Getenv("PORT"); httpPort == "" {
			log.Panicf("failed to read the PORT env variable to the application deployment")
		}
	}

	return fmt.Sprintf(":%s", httpPort)
}

// setupRSAKeys is the function that configures the RSA keys.
func setupRSAKeys() (authpkg.RSAKeys, error) {
	publicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to read the RSA public key file: %s", err.Error())
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return authpkg.RSAKeys{}, fmt.Errorf("failed to parse the RSA public key: %s", err.Error())
	}

	privateKey, err := ioutil.ReadFile(privateKeyPath)
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

	if deploy == "YES" {
		if dbURL = os.Getenv("DATABASE_URL"); dbURL == "" {
			return nil, fmt.Errorf("failed to read the DATABASE_URL environment variable to the application deployment")
		}
	}

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
