package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/weeback/grpc-project-template/internal/application/hello"
	"github.com/weeback/grpc-project-template/internal/auth"
	"github.com/weeback/grpc-project-template/internal/config"
	"github.com/weeback/grpc-project-template/internal/infrastructure/mongodb"
	"github.com/weeback/grpc-project-template/internal/infrastructure/transport/grpc"
	"github.com/weeback/grpc-project-template/internal/infrastructure/transport/http"
	"github.com/weeback/grpc-project-template/pkg"
	"github.com/weeback/grpc-project-template/pkg/net"

	hellopb "github.com/weeback/grpc-project-template/pb/hello"

	"github.com/gorilla/mux"

	googlegrpc "google.golang.org/grpc"
	googlealts "google.golang.org/grpc/credentials/alts"
	"google.golang.org/grpc/keepalive"
)

var (
	mongoURL = config.GetMongoURI()
	// captchaSecretKey = config.GetCloudflareTurnstileCredentials()

)

func init() {

	// Load the configuration from the JSON file; If it not error,
	// you can use config.GetOptionFirebaseAdmin function to get the Firebase admin options.
	if _, err := config.LoadWithJsonFile(config.GetFirebaseSdkCredentials()); err != nil {
		fmt.Printf("failed to load firebase admin options: %v\n", err)
		os.Exit(1)
	}

	// Load ALTS options
	if _, err := config.LoadALTS(); err != nil {
		fmt.Printf("failed to load ALTS options: %v\n", err)
		os.Exit(1)
	}

}

func main() {
	// Initialize the context
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

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

	auth := auth.New()

	// This is listed Google service accounts defined to allow accepting requests from Cloud Run.
	// If empty, it will allow every request.
	altsOpt := config.GetOptionALTS()
	expectedServiceAccounts := altsOpt.TargetServiceAccounts

	// gRPC servers can use ALTS credentials to allow clients to connect to them,
	// as illustrated next:
	grpcOpt := config.GetOptionGRPC()
	// Create gRPC server with increased timeouts and keepalive settings
	inst := googlegrpc.NewServer(
		googlegrpc.Creds(googlealts.NewServerCreds(googlealts.DefaultServerOptions())),
		googlegrpc.MaxRecvMsgSize(config.GetOptionGRPC().MaxRecvMsgSize),
		googlegrpc.MaxSendMsgSize(config.GetOptionGRPC().MaxSendMsgSize),
		googlegrpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcOpt.KeepaliveTime,
			Timeout: grpcOpt.KeepaliveTimeout,
		}),
		googlegrpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcOpt.KeepaliveTime,
			PermitWithoutStream: true,
		}),
		googlegrpc.ConnectionTimeout(grpcOpt.ConnectionTimeout),
		googlegrpc.StreamInterceptor(net.StreamInterceptor()),
		googlegrpc.UnaryInterceptor(net.UnaryServerAuthInterceptor(expectedServiceAccounts, auth.AuthFunc)),
	)
	//
	defer inst.Stop()
	//
	hellopb.RegisterHelloServiceServer(inst, grpc.NewHelloServiceHandler(helloRepo))

	/** Apply middleware to the router HTTP/1 (RESTful API)
	- Logging API request
	- Config CORS option (fix/access cors-domain problem)
	- Add middleware functions
	*/
	httpServer := net.Middleware(router, true)
	grpcServer := net.Walk(inst)
	mixed := net.MixHttp2(httpServer, grpcServer)

	// Open and listen port (:8080)
	if err := net.HttpServerWithConfig(":8080", mixed).ListenAndServe(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
