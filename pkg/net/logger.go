package net

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	// Header constants
	hostname = func() string {
		name, err := os.Hostname()
		if err != nil {
			name = "app"
		}
		// Append the process ID to the hostname for uniqueness
		return fmt.Sprintf("%s#%d", name, os.Getpid())
	}()
)

func apiLoggerHandler(h http.Handler) http.Handler {

	// Define the API logger
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// make a new response writer with the original writer
		wc := NewHttpWriter(w)

		defer printLogApi(wc, r, time.Now())

		// Call the next handler
		h.ServeHTTP(wc, r)
	})
}

func printLogApi(wc *ResponseWriter, r *http.Request, t time.Time) {
	var (
		origin    = r.Header.Get(headerOrigin)
		userAgent = r.Header.Get(headerUserAgent)
		requestId = r.Header.Get(xApiRequestId)
		clientId  = r.Header.Get(xApiClientId)
		more      string
	)
	// Add the request-Id to log
	if requestId == "" {
		requestId = "N/A"
	}
	if origin == "" {
		origin = "N/A"
	}
	if userAgent == "" {
		userAgent = "[N/A]"
	}
	if clientId == "" {
		clientId = "N/A"
	}
	//
	more = fmt.Sprintf("%s\t- - %s - %s - RequestID=%s - ClientID=%s", more, origin, userAgent, requestId, clientId)

	// Collect response headers
	cH := wc.Header().Clone()

	// Check if there is an error message
	if msg := cH.Get(xApiMoreError); msg != "" {
		more = fmt.Sprintf("%s > %s", more, msg)
	}

	// Log the request
	fmt.Printf("%s - %s | %20s --> %d %s - - %6s - %s - %s%s\n", time.Now().Format(time.DateTime), hostname,
		r.RemoteAddr, wc.StatusCode(), wc.Status(), r.Method, r.URL.String(), time.Since(t).String(), more)
}
