package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github/weeback/grpc-project-template/internal/application/hello"
	"github/weeback/grpc-project-template/internal/config"
	"github/weeback/grpc-project-template/internal/infrastructure/mongodb"
	"github/weeback/grpc-project-template/internal/infrastructure/transport/grpc"
	"github/weeback/grpc-project-template/internal/infrastructure/transport/http"
	"github/weeback/grpc-project-template/pkg"
	"github/weeback/grpc-project-template/pkg/net"

	// Import unuse package 'github/weeback/grpc-project-template/pkg' with tag '_' to execute init() function of global package.
	//	_ "github/weeback/grpc-project-template/pkg"

	pb "github/weeback/grpc-project-template/pb/hello"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	googlegrpc "google.golang.org/grpc"
	googlealts "google.golang.org/grpc/credentials/alts"
)

var (
	mongoURL = config.GetMongoURI()
	// captchaSecretKey = config.GetCloudflareTurnstileCredentials()

	logger *zap.Logger
)

func init() {
	//
	// Initialize the logger
	switch config.GetDeploymentEnvironment() {
	case config.Production:
		// Initialize logger
		prod, err := zap.NewProduction()
		if err != nil {
			logger = zap.NewExample()
			logger.Debug("failed to initialize production logger", zap.Error(err))
		} else {
			logger = prod
		}

	case config.Development:
		// Initialize logger
		dev, err := zap.NewDevelopment()
		if err != nil {
			logger = zap.NewExample()
			logger.Debug("failed to initialize production logger", zap.Error(err))
		} else {
			logger = dev
		}
	}

	// Load the configuration from the JSON file; If it not error,
	// you can use config.GetOptionFirebaseAdmin function to get the Firebase admin options.
	if _, err := config.LoadWithJsonFile(config.GetFirebaseSdkCredentials()); err != nil {
		fmt.Printf("failed to load firebase admin options: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	// Initialize the context
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)

	defer func(ctxCancelFunc context.CancelFunc, logger *zap.Logger) {
		ctxCancelFunc()
		if logger != nil {
			if err := logger.Sync(); err != nil {
				fmt.Printf("failed to sync logger: %v\n", err)
				os.Exit(1)
			}
		}
	}(cancel, logger)

	// Init connection
	databaseInter := mongodb.NewMongoDB(ctx, mongoURL)

	// captchaService := cloudflare.NewCaptchaService(captchaSecretKey, cloudflare.DefaultTurnstileVerifyURL)

	// Use mock-up service for testing
	// captchaService = mock.CaptchaService

	// =============================

	helloRepo := hello.NewHelloServiceRepo(databaseInter.ExampleDB)

	// Create a new router
	router := mux.NewRouter()

	// Register HTTP/1 handlers for your RESTful API service here
	httpHandler := http.NewHelloServiceHandler(helloRepo)
	// Register server path and handler
	router.HandleFunc("/healthcheck",
		pkg.HealthCheckHandler).Methods(http.MethodGet)

	if config.GetDeploymentEnvironment() == config.Development {
		// Register the mock handler for development environment
		router.PathPrefix("/proto/").Handler(
			net.FileServer("/proto/", "./proto/")).Methods(http.MethodGet)
	} else {
		router.PathPrefix("/proto/").Handler(
			pkg.PermissionDeniedHandler()).Methods(http.MethodGet)
	}

	// Register the SayHello handler
	router.HandleFunc("/say-hello", httpHandler.SayHello).Methods(http.MethodPost)
	// TODO: add more handlers here, template below:
	//
	// router.HandleFunc("/<path-to-entrypoint>", <handler-function-name>).Methods(<http-method(s)>)

	// This is listed Google service accounts defined to allow accepting requests from Cloud Run.
	// If empty, it will allow every request.
	expectedServiceAccounts := make([]string, 0)

	// gRPC servers can use ALTS credentials to allow clients to connect to them,
	// as illustrated next:
	inst := googlegrpc.NewServer(
		googlegrpc.Creds(googlealts.NewServerCreds(googlealts.DefaultServerOptions())),
		googlegrpc.MaxRecvMsgSize(config.GetOptionGRPC().MaxRecvMsgSize),
		googlegrpc.MaxSendMsgSize(config.GetOptionGRPC().MaxSendMsgSize),
	)
	//
	defer inst.Stop()
	//
	pb.RegisterHelloServiceServer(inst, grpc.NewHelloServiceHandler(helloRepo))

	/** Apply middleware to the router HTTP/1 (RESTful API)
	- Logging API request
	- Config CORS option (fix/access cors-domain problem)
	- Add middleware functions
	*/
	httpServer := net.Middleware(router, true)

	grpcServer := net.AllowServiceAccounts(inst, expectedServiceAccounts)

	mixed := net.MixHttp2(httpServer, grpcServer)

	/** Open and listen port (:8080)
	- The HttpServerWithConfig function is shortened from the following code:
	```
		def := &http.Server{
			Addr: "0.0.0.0:8080",
			Handler: mixed,
			ReadTimeout: 5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout: 15 * time.Second,
		}

		if err := def.ListenAndServe(); err != nil {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	```
	-It shortens and allows customization within the package without messing up the main function.
	*/
	if err := net.HttpServerWithConfig(":8080", mixed).ListenAndServe(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
