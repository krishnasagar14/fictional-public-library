package server

import (
	"fictional-public-library/config"
	"fictional-public-library/routerconfig"
	"github.com/gorilla/mux"
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
	r.initializeMiddleware(routerCfg)
	r.initializeRoutes(routerCfg.WebServerConfig)
}

// initializeRoutes ...
func (r *Router) initializeRoutes(webserverCfg *config.WebServerConfig) {}

// initializeMiddleware ..
func (r *Router) initializeMiddleware(routerCfg *routerconfig.RouterConfig) {}
