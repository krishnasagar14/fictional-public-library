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

	versionOne := "/v1"

	s.HandleFunc(versionOne+"/book/add", handlers.AddBookHandler()).
		Methods(http.MethodPost, http.MethodOptions).Name("AddBook")

	s.HandleFunc(versionOne+"/book/delete", handlers.DeleteBookHandler()).
		Methods(http.MethodDelete, http.MethodOptions).Name("DeleteBook")

	s.HandleFunc(versionOne+"/book/update", handlers.UpdateBookHandler()).
		Methods(http.MethodPut, http.MethodOptions).Name("UpdateBook")

	s.HandleFunc(versionOne+"/book/rent", handlers.RentBookHandler()).
		Methods(http.MethodPatch, http.MethodOptions).Name("RentBook")

	s.HandleFunc(versionOne+"/books", handlers.FetchAllBooksInLibraryHandler()).
		Methods(http.MethodGet, http.MethodOptions).Name("FetchAllBooks")
}
