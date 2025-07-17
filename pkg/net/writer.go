package net

import "net/http"

type ResponseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *ResponseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func (rw *ResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriter) StatusCode() int {
	return rw.status
}

func (rw *ResponseWriter) Status() string {
	return http.StatusText(rw.status)
}

func NewHttpWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w, status: 200}
}
