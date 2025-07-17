package net

import (
	"encoding/json"
	"io"
	"net/http"
)

var (
	firstOf = func(args []string) string {
		if len(args) > 0 {
			return args[0]
		}
		return ""
	}
)

func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

func QueryParams(r *http.Request, key string, def ...string) string {

	if val := r.URL.Query().Get(key); val != "" {
		return val
	}
	return firstOf(def)
}

func GetQueryParams(r *http.Request, key string) []string {
	return r.URL.Query()[key]
}

func GetRawData(r *http.Request) ([]byte, error) {

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			println(err.Error())
		}
	}(r.Body)

	return io.ReadAll(r.Body)
}

func ShouldBindJSON(r *http.Request, v interface{}) (raw []byte, err error) {

	raw, err = GetRawData(r)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(raw, v); err != nil {
		return raw, err
	}

	return raw, nil
}
