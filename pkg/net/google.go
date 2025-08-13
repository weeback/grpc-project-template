package net

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/weeback/grpc-project-template/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/alts"
	"google.golang.org/grpc/status"
)

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

// UnaryServerLoggingInterceptor creates a server interceptor for logging gRPC requests

// UnaryServerLoggingInterceptor creates a server interceptor for logging gRPC requests
func UnaryServerLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()

		// Extract request ID if available
		var reqID string
		if r, ok := req.(interface{ GetReqId() string }); ok && r != nil {
			reqID = r.GetReqId()
		}

		// Create logger with request context
		reqLogger := logger.NewEntry().With(
			zap.String("method", info.FullMethod),
			zap.String("req_id", reqID),
		)

		// Use the context with the logger
		ctx = logger.SetLoggerToContext(ctx, reqLogger)

		// Log the start of the request
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
	entry := logger.GetLoggerFromContext(ctx)
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
