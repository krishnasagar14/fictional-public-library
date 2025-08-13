package mw

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.RequestURI)
		next.ServeHTTP(resp, req)
	})
}
