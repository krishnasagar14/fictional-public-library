package tracing

import (
	"context"
	"fictional-public-library/literals"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		traceID := uuid.New().String()
		ctx := r.Context()
		ctx = context.WithValue(ctx, literals.TraceIDContextKey, traceID)
		ctx = context.WithValue(ctx, literals.EntrypointKey, mux.CurrentRoute(r).GetName())

		r = r.WithContext(ctx)

		w.Header().Set(literals.TraceIDContextKey, traceID)

		next.ServeHTTP(w, r)
	})
}
