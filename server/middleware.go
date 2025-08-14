package server

import (
	"context"
	"fictional-public-library/literals"
	"fictional-public-library/logging"
	"io/ioutil"
	"net/http"
)

const (
	jsonContentType = "application/contracts"
)

func AddContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("content-type", jsonContentType)

		// Token is valid, call the next handler
		next.ServeHTTP(w, r)
	})
}

func ReadReqMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				logging.Log.Error("Error while reading from request body: ", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, literals.RequestBody, data)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
