package net

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func Middleware(ro *mux.Router, enableCORS bool, middlewareFunc ...http.HandlerFunc) http.Handler {
	// Define CORS options
	corsHandler := func(h http.Handler) http.Handler {
		if !enableCORS {
			return h
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(corsAllowOriginHeader, "*")
			w.Header().Set(corsAllowMethodsHeader, "GET, POST, OPTIONS")
			w.Header().Set(corsAllowHeadersHeader, "*")
			w.Header().Set(corsMaxAgeHeader, "3600")
			w.Header().Set(corsAllowCredentialsHeader, "true")
			h.ServeHTTP(w, r)
		})
	}
	// Apply the middleware functions
	middlewareHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set the request-Id
			if requestId := r.Header.Get(xApiRequestId); requestId == "" {
				r.Header.Set(xApiRequestId, uuid.NewString())
			}
			for _, Func := range middlewareFunc {
				Func(w, r)
			}
			h.ServeHTTP(w, r)
		})
	}

	//
	ro.Use(corsHandler, middlewareHandler)

	// Walk through all the registered routes
	err := ro.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			fmt.Printf("[%-8s] %s\n", "", pathTemplate)
			return nil
		}
		for _, method := range methods {
			fmt.Printf("[%-8s] %s\n", method, pathTemplate)
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking routes: ", err)
	}

	return apiLoggerHandler(ro)
}
