package server

import (
	"fictional-public-library/handlers"
	"fictional-public-library/routerconfig"
	"fictional-public-library/tracing"
	"github.com/gorilla/mux"
	"net/http"
)

// Router ...
type Router struct {
	*mux.Router
}

// NewRouter ...
func NewRouter() *Router {
	return &Router{mux.NewRouter()}
}

// InitializeRouter ...
func (r *Router) InitializeRouter(routerCfg *routerconfig.RouterConfig) {
	// middlewares
	r.Use(AddContentTypeMiddleware)
	r.Use(tracing.TraceMiddleware)
	r.Use(ReadReqMiddleware)

	s := (*r).PathPrefix(routerCfg.WebServerConfig.RoutePrefix).Subrouter()

	s.HandleFunc("/book/add", handlers.AddBookHandler(routerCfg)).Methods("POST").
		Methods(http.MethodPost, http.MethodOptions).Name("AddBook")
}
