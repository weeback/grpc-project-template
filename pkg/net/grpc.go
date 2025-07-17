package net

import (
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

func MixHttp2(rest, gRPC http.Handler) http.Handler {
	// Trộn cả gRPC và REST mux
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get(headerContentType), "application/grpc") {
			gRPC.ServeHTTP(w, r)
		} else {
			rest.ServeHTTP(w, r)
		}
	})

	return h2c.NewHandler(mainHandler, &http2.Server{})
}

// Walk for gRPC only
func Walk(inst *grpc.Server, expectedServiceAccounts ...string) http.Handler {
	// Print gRPC service information
	for key, inf := range inst.GetServiceInfo() {
		for _, mt := range inf.Methods {
			mode := "Simple RPC"
			if mt.IsServerStream {
				if mt.IsClientStream {
					mode = "Bidirectional streaming RPC"
				} else {
					mode = "Server-side streaming RPC"
				}
			} else if mt.IsClientStream {
				mode = "Client-side streaming RPC"
			}
			txt := fmt.Sprintf("%-30s | %s.(%s)", mode, key, mt.Name)
			println("[gRPC]", txt)
		}
	}
	return AllowServiceAccounts(inst, expectedServiceAccounts)
}
