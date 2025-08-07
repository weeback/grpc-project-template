package http

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"net/http"

	"github.com/weeback/grpc-project-template/internal/entity/hello"
	"github.com/weeback/grpc-project-template/pkg"
	"github.com/weeback/grpc-project-template/pkg/jwt"
	"github.com/weeback/grpc-project-template/pkg/net"

	common "github.com/weeback/grpc-project-template/pb/common"
	hellopb "github.com/weeback/grpc-project-template/pb/hello"
)

func NewHelloServiceHandler(svc hello.Repository) *HelloServiceHandler {
	return &HelloServiceHandler{
		service: svc,
		validateAuthenticationFunc: func(ctx context.Context, c *common.ClientJwt, v proto.Message) (context.Context, error) {

			// TODO: Implement your Public key get from user logic here
			// For example, you might want to fetch the public key from a database or a configuration file.
			// This is just a placeholder implementation.
			// You can replace this with your actual public key retrieval logic.
			keyPair, _ := jwt.GenerateKeyPair()

			claims, err := jwt.ParseClaimsWithoutVerification(keyPair.PublicKey, c.Jwt)
			if err != nil {
				// Placeholder for JWT validation logic
				// You can implement your JWT validation logic here

				// For example, you might want to check the signature, expiration, etc.
				// If the JWT is valid, return nil; otherwise, return an error.
				// This is just a placeholder implementation.
			}
			if claims == nil {
				return ctx, fmt.Errorf("invalid JWT token")
			}
			if err := claims.ParsePayload(v); err != nil {
				return ctx, status.Errorf(codes.Unauthenticated, "invalid JWT token, error %#v", err)
			}
			// If the JWT is valid, return the claims and apply them to the context
			return claims.ApplyContext(ctx, c.ReqId), nil
		},
	}
}

type HelloServiceHandler struct {
	service                    hello.Repository
	validateAuthenticationFunc func(ctx context.Context, c *common.ClientJwt, v proto.Message) (context.Context, error)
}

func (h *HelloServiceHandler) SayHello(w http.ResponseWriter, r *http.Request) {
	var (
		request hellopb.HelloRequest
	)
	// Read the request body from http
	if raw, err := net.ShouldBindJSON(r, &request); err != nil {
		net.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("failed to parse request: %v.\r\n%s", err, string(raw)))
		return
	}
	// Redirect sends request to login service
	resp, err := h.service.SayHello(r.Context(), &request)
	if err != nil {
		fmt.Printf("api-%s :%s\n", r.RequestURI, err.Error())
	}
	// Write response to http
	if err := net.WriteJSONbyError(w, http.StatusOK, err, resp); err != nil {
		net.WriteError(w, http.StatusServiceUnavailable, err)
	}
}

func (h *HelloServiceHandler) UseStandardResponse(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		in  common.ClientJwt

		request hellopb.PayloadRequest
	)
	// Read the request body from http
	if raw, err := net.ShouldBindJSON(r, &in); err != nil {
		net.WriteError(w, http.StatusBadRequest,
			fmt.Errorf("failed to parse request: %v.\r\n%s", err, string(raw)))
		return
	}

	// TODO: you can add some code here
	// to handle before forward call service
	// Validate the JWT token
	// This is a placeholder for JWT validation logic
	// You should implement your JWT validation logic here
	// For example, you might want to check the signature, expiration, etc.
	// Assuming h.validateAuthenticationFunc() is a method that checks the validity of the JWT token
	// If h.validateAuthenticationFunc() is not defined, you can implement your own validation logic
	// or remove this comment if not needed.
	jwtCtx, err := h.validateAuthenticationFunc(ctx, &in, &request)
	if err != nil {
		net.WriteError(w, http.StatusUnauthorized,
			fmt.Errorf("failed to validate JWT token: %v", err))
		return
	}
	// Redirect sends request to login service
	resp, err := h.service.UseStandardResponse(jwtCtx, &request)
	if err != nil {
		fmt.Printf("api-%s :%s\n", r.RequestURI, err.Error())
	}
	// Write response to http
	if err := net.WriteJSONbyError(w, http.StatusOK, err, resp); err != nil {
		net.WriteError(w, http.StatusServiceUnavailable, err)
	}
}

func (h *HelloServiceHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {

	version := pkg.GetBuiltVersionInfo()

	/** Write headers
	- Content-Type: text/plain
	- X-Content-Type-Options: nosniff
	https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options
	*/
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	/** Write status code 200
	(https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/200)
	*/
	w.WriteHeader(http.StatusOK)

	/** Write data for response, example:
	>> Version: V1.0.0
	>> Build by Admin
	>> Build at 2025-04-12T00:00:00Z
	*/
	if _, warn := bytes.NewBufferString(version).WriteTo(w); warn != nil {
		println(warn.Error())
	}
}
