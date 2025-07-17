package net

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetWd() string {
	wdPath, err := os.Getwd()
	if err != nil {
		log.Printf("can not get path of current directory: %s", err.Error())
		return "."
	}
	return wdPath
}

func FileServer(prefix, dirPath string) http.Handler {
	fs := http.FileServer(http.Dir(dirPath))
	return http.StripPrefix(prefix, fs)
}

func HttpServerWithConfig(addr string, handler http.Handler) *http.Server {
	//
	def := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	//
	defer func() {
		log.Printf("Working directory: %s%s\n", GetWd(), string(filepath.Separator))
		log.Printf("Serving on %s\n", def.Addr)
	}()
	//
	if addr != "" && addr != def.Addr {
		if strings.HasPrefix(addr, ":") {
			def.Addr = "0.0.0.0" + addr
		} else {
			def.Addr = addr
		}
	}
	return def
}
