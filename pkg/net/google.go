package net

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/alts"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// AllowServiceAccounts creates a middleware that checks if the request is authorized by service accounts
// and allows gRPC requests to pass through if they are valid.
//
// Deprecated: Use `UnaryServerAuthInterceptor` for gRPC server-side authentication.
func AllowServiceAccounts(inst *grpc.Server, expectedServiceAccounts []string) http.Handler {
	if len(expectedServiceAccounts) == 0 {
		return inst
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Chỉ kiểm tra ClientAuthorizationCheck nếu request là gRPC (HTTP/2)
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get(headerContentType), "application/grpc") {
			if serviceAccount := r.Header.Get(xApiServiceAccount); serviceAccount != "" {
				// Header `X-Service-Account` ready -> case test by Postman or a script
				if !slices.Contains(expectedServiceAccounts, serviceAccount) {
					// write status
					w.WriteHeader(http.StatusUnauthorized)
					// write a message
					if _, err := w.Write([]byte("Unauthorized: Invalid service account")); err != nil {
						return
					}
					return
				}
			} else {
				// Default case Cloud Run Service
				if err := clientAuthorizationCheck(r.Context(), expectedServiceAccounts); err != nil {
					// write status
					w.WriteHeader(http.StatusUnauthorized)
					// write a message
					if _, err := w.Write([]byte("Unauthorized: Invalid service account")); err != nil {
						return
					}
					return
				}
			}
		}
		// If it is REST or a gRCC request valid, forward to inst.ServeHTTP
		inst.ServeHTTP(w, r)
	})
}

func StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		fmt.Printf("Stream started: %v\n", info.FullMethod)
		if err := handler(srv, ss); err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return err
		}
		fmt.Printf("Stream completed successfully\n")
		return nil
	}
}

// UnaryServerAuthInterceptor creates a server interceptor for attack middleware function to gRPC requests
func UnaryServerAuthInterceptor(expectedServiceAccounts []string, authFunc func(fullMethod string, bodyHash string, jwtStr string) error) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var (
			startTime = time.Now()
			reqID     string
			md        metadata.MD

			bodyHash   string
			jwtAuthStr string
		)

		//  Extract Metadata from context
		if fromCtx, ok := metadata.FromIncomingContext(ctx); ok && fromCtx != nil {
			md = fromCtx
		}
		if md != nil {
			// Extract authorization from metadata
			if auth := md.Get(headerAuthorization); len(auth) > 0 {
				if strings.HasPrefix(auth[0], "Bearer ") {
					jwtAuthStr = strings.TrimPrefix(auth[0], "Bearer ")
				} else {
					jwtAuthStr = auth[0]
				}
			}
			// Extract request ID from metadata
			if reqIDs := md.Get(xApiRequestId); len(reqIDs) > 0 {
				reqID = reqIDs[0]
			}
		}

		// Extract request ID if available
		if r, ok := req.(interface{ GetReqId() string }); ok && r != nil {
			reqID = r.GetReqId()
		}

		// Create logger with request context
		reqLogger := getLogEntry().With(
			zap.Bool("proto_marshaled", false),
			zap.String("method", info.FullMethod),
			zap.String("req_id", reqID),
			zap.Any("metadata", md),
		)
		// Use the context with the logger
		ctx = setLoggerToContext(ctx, reqLogger)

		// Check if the request is authorized by service account
		if err := clientAuthorizationCheck(ctx, expectedServiceAccounts); err != nil {
			// Log the error
			reqLogger.Error("Client authorization check failed",
				zap.String("method", info.FullMethod),
				zap.String("req_id", reqID),
				zap.Error(err),
			)
			return nil, err
		}

		msg, ok := req.(proto.Message)
		if ok {
			// Marshal the proto message to log its SHA256 hash
			b, err := proto.Marshal(msg)
			if err != nil {
				reqLogger = reqLogger.With(
					zap.Any("request", msg), // If the request is a proto message, log it
					zap.Errors("marshal_error", []error{err}),
				)
			} else {
				sum := sha256.Sum256(b)
				bodyHash = hex.EncodeToString(sum[:])
				reqLogger = reqLogger.With(
					zap.Bool("proto_marshaled", true),
					zap.String("sum", hex.EncodeToString(sum[:])),
					zap.String("request", string(b)), // If the request is a proto message, log it
				)
			}
		} else {
			// Otherwise, log the request as a generic interface
			reqLogger = reqLogger.With(
				zap.Any("request", req),
				zap.Errors("proto_marshal_error", []error{status.Errorf(codes.Internal, "request is not a proto message")}),
			)
		}

		if err := authFunc(info.FullMethod, bodyHash, jwtAuthStr); err != nil {
			reqLogger.Error("Authorization failed",
				zap.String("body_hash", bodyHash),
				zap.String("jwt", jwtAuthStr),
				zap.Error(err))
			return nil, status.Errorf(codes.Unauthenticated, "Authorization failed: %v", err)
		}

		//
		reqLogger.Debug("gRPC middleware interceptor", zap.Duration("duration", time.Since(startTime)))

		reqLogger.Info("gRPC request started", zap.Time("start_time", startTime))
		// Process the request
		resp, err := handler(ctx, req)

		// Get status code
		statusCode := codes.OK
		if err != nil {
			statusCode = status.Code(err)
		}

		// Log completion
		reqLogger.Info("gRPC request completed",
			zap.Any("response", resp),
			zap.String("status", statusCode.String()),
			zap.Duration("duration", time.Since(startTime)),
			zap.Error(err),
		)

		return resp, err
	}
}

// UnaryServerLoggingInterceptor creates a server interceptor for logging gRPC requests
func UnaryServerLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		// Extract request ID if available
		var reqID string
		if r, ok := req.(interface{ GetReqId() string }); ok && r != nil {
			reqID = r.GetReqId()
		}

		// Create logger with request context
		reqLogger := getLogEntry().With(
			zap.Bool("proto_marshaled", false),
			zap.String("method", info.FullMethod),
			zap.String("req_id", reqID),
		)

		// Use the context with the logger
		ctx = setLoggerToContext(ctx, reqLogger)

		// Log the start of the request
		startTime := time.Now()
		reqLogger.Info("gRPC request started",
			zap.Any("request", req),
			zap.Time("start_time", startTime),
		)

		// Process the request
		resp, err := handler(ctx, req)

		// Get status code
		statusCode := codes.OK
		if err != nil {
			statusCode = status.Code(err)
		}

		// Log completion
		reqLogger.Info("gRPC request completed",
			zap.Any("response", resp),
			zap.String("status", statusCode.String()),
			zap.Duration("duration", time.Since(startTime)),
			zap.Error(err),
		)

		return resp, err
	}
}

func clientAuthorizationCheck(ctx context.Context, expectedServiceAccounts []string) error {
	if len(expectedServiceAccounts) == 0 {
		return nil // No service accounts to check against, allow all
	}
	authInfo, err := alts.AuthInfoFromContext(ctx)
	if err != nil {
		return status.Errorf(codes.PermissionDenied, "The context is not an ALTS-compatible context: %v", err)
	}
	entry := getLoggerFromContext(ctx)
	entry.Debug("ALTS AuthInfo",
		zap.String("PeerServiceAccount", authInfo.PeerServiceAccount()),
		zap.String("LocalServiceAccount", authInfo.LocalServiceAccount()),
		zap.String("ApplicationProtocol", authInfo.ApplicationProtocol()),
		zap.String("RecordProtocol", authInfo.RecordProtocol()))

	peer := authInfo.PeerServiceAccount()
	entry.Debug("ALTS AuthInfo", zap.String("peer", peer))

	for _, sa := range expectedServiceAccounts {
		if strings.EqualFold(peer, sa) {
			return nil
		}
	}
	return status.Errorf(codes.PermissionDenied, "Client %v is not authorized", peer)
}
