package grpc

import (
	"context"

	"github.com/weeback/grpc-project-template/internal/entity/hello"
	"github.com/weeback/grpc-project-template/pkg/jwt"

	common "github.com/weeback/grpc-project-template/pb/common"
	hellopb "github.com/weeback/grpc-project-template/pb/hello"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewHelloServiceHandler creates a new HelloServiceHandler
func NewHelloServiceHandler(svc hello.Repository) hellopb.HelloServiceServer {
	return &HelloServiceHandler{
		service: svc,
		validateAuthenticationFunc: func(ctx context.Context, c *common.ClientJwt) (*jwt.MapClaims, context.Context, error) {

			// TODO: Implement your Public key get from user logic here
			// For example, you might want to fetch the public key from a database or a configuration file.
			// This is just a placeholder implementation.
			// You can replace this with your actual public key retrieval logic.
			keyPair, _ := jwt.GenerateKeyPair()

			claims, err := jwt.ParseClaims(keyPair.PublicKey, c.Jwt)
			if err != nil {
				// Placeholder for JWT validation logic
				// You can implement your JWT validation logic here

				// For example, you might want to check the signature, expiration, etc.
				// If the JWT is valid, return nil; otherwise, return an error.
				// This is just a placeholder implementation.
			}
			if claims == nil {
				return nil, ctx, status.Errorf(codes.Unauthenticated, "invalid JWT token")
			}
			// If the JWT is valid, return the claims and apply them to the context
			return claims, claims.ApplyContext(ctx, c.ReqId), nil
		},
	}
}

type HelloServiceHandler struct {
	hellopb.UnimplementedHelloServiceServer
	service                    hello.Repository
	validateAuthenticationFunc func(ctx context.Context, c *common.ClientJwt) (*jwt.MapClaims, context.Context, error)
}

func (h *HelloServiceHandler) SayHello(ctx context.Context, request *hellopb.HelloRequest) (*hellopb.HelloReply, error) {
	// TODO: you can add someone code here
	// to handle before forward call service
	return h.service.SayHello(ctx, request)
}

func (h *HelloServiceHandler) UseStandardResponse(ctx context.Context, in *common.ClientJwt) (*common.StandardResponse, error) {
	// TODO: you can add some code here
	// to handle before forward call service
	// Validate the JWT token
	// This is a placeholder for JWT validation logic
	// You should implement your JWT validation logic here
	// For example, you might want to check the signature, expiration, etc.
	// Assuming h.validateAuthenticationFunc() is a method that checks the validity of the JWT token
	// If h.validateAuthenticationFunc() is not defined, you can implement your own validation logic
	// or remove this comment if not needed.
	claims, jwtCtx, err := h.validateAuthenticationFunc(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "failed to validate JWT token: %v", err)
	}
	// Parse the payload to the expected type using type detection
	request, ok := claims.Payload.(*hellopb.PayloadRequest)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payload type in JWT claims")
	}

	// Apply the claims to the context, and forward the request to the service
	return h.service.UseStandardResponse(jwtCtx, request)
}
