package commands

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/icaroribeiro/new-go-code-challenge-template/docs/api/swagger"
	healthcheckhandler "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/http/presentation/handler/healthcheck"
	healthcheckrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/http/router/healthcheck"
	swaggerrouter "github.com/icaroribeiro/new-go-code-challenge-template/internal/transport/http/router/swagger"
	envpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/env"
	handlerhttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/handler"
	routehttputilpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/httputil/route"
	serverpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/server"
	"github.com/spf13/cobra"
	//validatorpkg "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator"
	//uuidvalidator "github.com/icaroribeiro/new-go-code-challenge-template/pkg/validator/uuid"
	//validatorv2 "gopkg.in/validator.v2"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the API",
	Run:   execRunCmd,
}

var (
	deploy = envpkg.GetEnvWithDefaultValue("DEPLOY", "NO")

	httpHost = envpkg.GetEnvWithDefaultValue("HTTP_HOST", "localhost")
	httpPort = envpkg.GetEnvWithDefaultValue("HTTP_PORT", "8080")
)

func execRunCmd(cmd *cobra.Command, args []string) {
	tcpAddress := setupTcpAddress()

	// validationFuncs := map[string]validatorv2.ValidationFunc{
	// 	"uuid": uuidvalidator.Validate,
	// }

	// validator, err := validatorpkg.New(validationFuncs)
	// if err != nil {
	// 	log.Panic(err.Error())
	// }

	routes := make(routehttputilpkg.Routes, 0)

	routes = append(routes, swaggerrouter.ConfigureRoutes()...)

	healthCheckHandler := healthcheckhandler.New()
	routes = append(routes, healthcheckrouter.ConfigureRoutes(healthCheckHandler)...)

	// fileHddStorageRepository := filehddstoragerepository.New(storage)
	// fileService := fileservice.New(fileHddStorageRepository, validator)
	// fileHandler := filehandler.New(fileService)
	// routes = append(routes, filerouter.ConfigureRoutes(fileHandler)...)

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

		return fmt.Sprintf(":%s", httpPort)
	}

	return fmt.Sprintf("%s:%s", httpHost, httpPort)
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
