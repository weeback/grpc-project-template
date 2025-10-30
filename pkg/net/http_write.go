package net

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/textproto"
)

// WriteJSON writes the JSON representation of v to the http.ResponseWriter.
// Payload body is JSON encoded of (v), and the Content-Type is set to application/json.
func WriteJSON(w http.ResponseWriter, httpStatus int, v interface{}) error {

	// validate http status code
	if http.StatusText(httpStatus) == "" {
		return fmt.Errorf("invalid http status code: %d", httpStatus)
	}
	// Set the Content-Type to application/json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// write the response
	w.WriteHeader(httpStatus)
	return json.NewEncoder(w).Encode(v)
}

// WriteJSONbyError writes the JSON representation of v to the http.ResponseWriter.
// Payload body is JSON encoded of (v) if err is nil, otherwise it returns the JSON encoded of error message and (v) is ignored.
func WriteJSONbyError(w http.ResponseWriter, httpStatus int, err error, v interface{}) error {

	// validate http status code
	if http.StatusText(httpStatus) == "" {
		return fmt.Errorf("invalid http status code: %d", httpStatus)
	}

	// Set the Content-Type to application/json
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// add more error message to the header
	if err != nil {
		w.Header().Set(textproto.CanonicalMIMEHeaderKey("x-more-error"), fmt.Sprintf("%+v", err))
	}
	// write the response
	w.WriteHeader(httpStatus)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, httpStatus int, err error) {

	write := func() error {
		// validate http status code
		if http.StatusText(httpStatus) == "" {
			return fmt.Errorf("invalid http status code: %d", httpStatus)
		}
		// Set the Content-Type to application/json
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// write the response
		w.WriteHeader(httpStatus)
		return json.NewEncoder(w).Encode(map[string]interface{}{
			"isError": true,
			"code":    111,
			"message": err.Error(),
		})
	}

	if warn := write(); warn != nil {
		println("\n" + warn.Error())
	}
}

func WriteBytes(w http.ResponseWriter, httpStatus int, data []byte) error {
	// validate http status code
	if http.StatusText(httpStatus) == "" {
		return fmt.Errorf("invalid http status code: %d", httpStatus)
	}
	// write the response
	w.WriteHeader(httpStatus)
	_, err := w.Write(data)
	return err
}
