###

### Required
(https://grpc.io/docs/languages/go/quickstart/)

1. You can install the protocol compiler, protoc, with a package manager under Linux, macOS, or Windows using the following commands.

   - Linux, using apt or apt-get, for example:
   ```
    $ apt install -y protobuf-compiler
    $ protoc --version  # Ensure compiler version is 3+
   ```

   - MacOS, using Homebrew:
   ```
    $ brew install protobuf
    $ protoc --version  # Ensure compiler version is 3+
   ```

   - Windows, using Winget
   ```
    $ winget install protobuf
    $ protoc --version # Ensure compiler version is 3+
   ```

2. Install the protocol compiler plugins for Go using the following commands:
   ```
    $ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    $ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

3. Update your PATH so that the protoc compiler can find the plugins:
   ```
    $ export PATH="$PATH:$(go env GOPATH)/bin"
   ```

---

### How to implement gRPC proto file, mix HTTP/1 and HTTP/2

1. Create a protobuf file [proto/hello.proto](proto/hello.proto), example:
   ```protobuf
   # filename=proto/hello.proto
   
   syntax = "proto3";

   option go_package = "github.com/weeback/grpc-project-template/pb";
   
   package pb;
   
   service HelloService {
    rpc SayHello (HelloRequest) returns (HelloReply);
   }
   
   message HelloRequest {
    string name = 1;
   }
   
   message HelloReply {
    string message = 1;
   }
   ```

2. Generate some necessary structure files from protobuf
   - Open `Terminal` and generate structure file:
   ```shell
   protoc --proto_path="proto" \
      --go_out="pb" --go_opt=paths=source_relative \
      --go-grpc_out="pb" --go-grpc_opt=paths=source_relative \
      proto/hello.proto;
   ```
   After execution, get the necessary structure files:
   ```
   ├── pb
   │   ├── hello.pb.go
   │   └── hello_grpc.pb.go
   └── proto
       └── hello.proto
   ```

   ___Create this code in a Makefile for reuse:___
   ```makefile

   # ====================================================================================
   # gRPC generate code
   # ====================================================================================

   ### --- existing script --- ###
   
   # Path Generating gRPC reflection code
   PROTO_PATH = proto
   PROTO_GEN = pb
   PROTO_PATH_HELLO = $(PROTO_PATH)

   # Gen target: Hello Service
   grpc-generate-hello: grpc-force-hello
    @echo "==> Generating gRPC code...";
    @protoc --proto_path=$(PROTO_PATH_HELLO) \
        --go_out=$(PROTO_GEN) --go_opt=paths=source_relative \
        --go-grpc_out=$(PROTO_GEN) --go-grpc_opt=paths=source_relative \
        $(PROTO_PATH_HELLO)/hello.proto || { echo "Failed to generate gRPC code"; exit 1; }
    @echo "==> gRPC code generated successful: $(PROTO_PATH_HELLO)/hello.pb.go, $(PROTO_PATH_HELLO)/hello_grpc.pb.go"

   # Force target
   grpc-force-hello:
    @echo "==> Creating output directory..."
    @mkdir -p $(PROTO_GEN)
    @echo "==> Cleaning up old generated files..."
    @rm -rf $(PROTO_GEN)/hello*.go

   ### --- existing script --- ###
   ```
   
3. Implement gRPC code

   * Create an [internal/entity/hello/repository.go](internal/entity/hello/repository.go) file to define
   the mapping entities/repositories for the methods defined by protobuf. 
   ```go
   /** filename=internal/entity/hello/repository.go
   */
   package hello
   
   type Repository interface {
    SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error)
   }
   ```
   
   * Create a file [internal/application/hello/controller.go](internal/application/hello/controller.go) to implement the logic for the application's functionality.
   You can define another file here `internal/application/hello/validation.go`, `internal/application/hello/utils.go` ...
   ```go
   /** filename=internal/application/hello/controller.go
   */
   package hello

   import (
       "context"
   
       "github.com/weeback/grpc-project-template/internal/entity/hello"
       "github.com/weeback/grpc-project-template/pb"
   )
   
   func NewHelloServiceRepo() hello.Repository {
    return &controller{
        // ...
    }
   }
   
   type controller struct {
    // TODO: define some fields here
   }
   
   func (ins *controller) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
    // TODO: implement logic me
   }
   ```
   
   * Create a file [internal/infrastructure/transport/grpc/hello_service.go](internal/infrastructure/transport/grpc/hello_service.go) to transfer data protobuf
   for logic service
   ```go
   /** filename=internal/infrastructure/transport/grpc/hello_service.go
   */
   package grpc

   import (
   "context"
   
       "github.com/weeback/grpc-project-template/internal/entity/hello"
       "github.com/weeback/grpc-project-template/pb"
   )
   
   func NewHelloServiceHandler(svc hello.Repository) *HelloServiceHandler {
    return &HelloServiceHandler{
        service: svc,
    }
   }
   
   type HelloServiceHandler struct {
    pb.HelloServiceServer
    service hello.Repository
   }
   
   func (h *HelloServiceHandler) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloReply, error) {
    // TODO: you can add someone code here
    // to handle before forward call service
    return h.service.SayHello(ctx, request)
   }
   ```
   
   ***==> Files mapping example:***
   
   ```
   ├── internal
   │   ├── entity
   │   │   └── hello
   │   │       └── repository.go
   │   └── infrastructure
   │       └── transport
   │           ├── grpc
   │           │   └── hello_service.go
   │           └── http
   ├── pb
   │   ├── common
   │   │   └── standard.pb.go
   │   └── wls
   │       ├── web_login_session.pb.go
   │       └── web_login_session_grpc.pb.go
   └── proto
       ├── common
       │   └── standard.proto
       └── wls
           └── web_login_session.proto
   ```

4. Deploy gRPC api server
   
   * In file `main.go`:
     - f
   ```go
   // filename=cmd/HelloService/main.go
   package main
   
   import (
      
      // === existing code ===
   
      "github.com/weeback/grpc-project-template/internal/application/hello"
      "github.com/weeback/grpc-project-template/internal/infrastructure/transport/grpc"
      "github.com/weeback/grpc-project-template/pb"
   
      googlegrpc "google.golang.org/grpc"
   )
   
   // === existing code ===
      inst := googlegrpc.NewServer()
      defer inst.Stop()
      
      helloRepo := hello.NewHelloServiceRepo()
      pb.RegisterHelloServiceServer(inst, grpc.NewHelloServiceHandler(helloRepo))
   
   // === existing code ===
   ```

5. Build a custom handler for the REST API (HTTP/1) \
   ___to configure mix run mode for two modes sharing one port 
   (because Google Cloud Run does not allow publishing multiple ports)___


---

### How to Deploy an application to Google Cloud Run?

#### Required:
1. Required roles

   To get the permissions that you need to deploy Cloud Run services, ask your administrator to grant you the following IAM roles:

   - Cloud Run Developer (roles/run.developer) on the Cloud Run service
   - Service Account User (roles/iam.serviceAccountUser) on the service identity
   - Artifact Registry Reader (roles/artifactregistry.reader) on the Artifact Registry repository of the deployed container image

   For a list of IAM roles and permissions that are associated with Cloud Run, see Cloud Run IAM roles and Cloud Run IAM permissions.
   If your Cloud Run service interfaces with Google Cloud APIs, such as Cloud Client Libraries, see the service identity configuration guide.
   For more information about granting roles, see deployment permissions and manage access.

   https://cloud.google.com/run/docs/deploying?authuser=1#required_roles

2. Login with a Google account and enable Google Cloud services

   - Authenticate with Google Cloud: Run the following command to authenticate:
   ```
    $ gcloud auth login
   ```

   - Set the Correct Project: Ensure the correct project is set:
   ```
    $ gcloud config set project service-test-40fef
   ```

   - Enable Artifact Registry API: Verify that the Artifact Registry API is enabled for your project:
   ```
    $ gcloud services enable artifactregistry.googleapis.com
   ```

   - Re-authenticate Docker with GCR: Reconfigure Docker authentication for the specified region:
   ```
    $ gcloud auth configure-docker asia.gcr.io
   ```

#### Build docker package and deploy 

1. Insert the target `build` code into the [Makefile](Makefile) file.
 
   - Use the `?=` operator to declare a variable, the purpose is to override the value passed in from
   the terminal command.
   - Insert the variable `MAIN_GO_FILE ?= <default-value` to customize the main.go file path, the purpose is
   to be able to build different services with different paths without having to change or create another target.
   - Insert the variable `BINARY ?= <default-value>` to customize the binary output file name, in case 
   the file name between different services is inconsistent.
   
   ```makefile
   # === existing script === #
   MAIN_GO_FILE ?= cmd/v1/*.go
   BINARY ?= bin/app
   # === existing script === #
   # Build application binary
   build:
      @echo "==> Building application..."
      @mkdir -p $(BIN_DIR)
      @export CGO_ENABLED=0; \
         go build -v $(LDFLAGS) -o $(BINARY) -trimpath $(MAIN_GO_FILE)
      @echo "==> Build successful: $(BINARY)"
   # === existing script === #
   ```
   
2. Create a file [[Name]Service.Dockerfile](HelloService.Dockerfile) 

   Change or adjust some content below as appropriate.

   ```dockerfile
   FROM amd64/golang:1.24-alpine3.21 AS build
   
   # === existing script === #

   RUN MAIN_GO_FILE="cmd/HelloService/*.go" BINARY="bin/HelloService" \
   make --makefile=Makefile build
   
   # Use the latest alpine image
   FROM alpine:3.21
   # === existing script === #
   # Copy the built Go application from the build stage
   COPY --from=build /app/bin/HelloService /app/bin/HelloService
   # === existing script === #
   CMD [ "/app/bin/HelloService" ]
   ```
   
3. Create a file [Makefile.CloudRun-[Name]Service](Makefile.CloudRun-HelloService)

   Note: Point the `docker build -f <filename>` command to the service's corresponding `HelloService.Dockerfile` file.

   ```makefile
   
   # === existing script === #
   build-docker: active-service-account
      @echo "Building docker image..."
      @docker build -t $(GOOGLE_CLOUD_IMAGE_NAME):$(GOOGLE_CLOUD_IMAGE_TAG) -f HelloService.Dockerfile . || { echo "Failed to build docker image"; exit 1; }
      @echo "Docker image built successfully"
   
   # === existing script === #
   ```
   
4. Deploy to Google Cloud Run

   - Build docker image, push to the artifact registry, and deploy
   ```shell
    $ make --makefile=Makefile.CloudRun-[NAME]Service deploy
   ```

   - If incorrect with error `ermission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock: Get "http://%2Fvar%2Frun%2Fdocker.sock/v1.47/containers/json": dial unix /var/run/docker.sock: connect: permission denied` to add user for group 'docker' and try again:
   ```shell
    $ sudo usermod -aG docker $USER
    $ newgrp docker
   ```