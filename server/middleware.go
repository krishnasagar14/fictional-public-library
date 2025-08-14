package server

import "net/http"

const (
	protobufContentType = "application/x-protobuf"
)

func AddContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("content-type", protobufContentType)

		// Token is valid, call the next handler
		next.ServeHTTP(w, r)
	})
}
