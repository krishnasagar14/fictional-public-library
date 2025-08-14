package server

import (
	"fictional-public-library/routerconfig"
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

	s := (*r).PathPrefix(routerCfg.WebServerConfig.RoutePrefix).Subrouter()

	s.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).
		Methods(http.MethodGet, http.MethodOptions).Name("public-library")
}
