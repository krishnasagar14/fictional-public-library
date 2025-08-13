package server

import (
	"fictional-public-library/config"
	mw "fictional-public-library/middlewares"
	main_routes "fictional-public-library/routers"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"fictional-public-library/db"
)

type Server struct {
	Configuration *config.WebServerConfig
	Router        *Router
}

// NewServer ...
func NewServer(config *config.WebServerConfig) *Server {
	server := &Server{
		Configuration: config,
		Router:        NewRouter(),
	}

	return server
}

func RunServer() (err error) {
	return nil
}

func main() {
	err := db.ConnectDB("local_db")
	if err != nil {
		os.Exit(1)
	}

	portNo := 9000
	fmt.Println("Starting server on port:", portNo)
	router := main_routes.RegisterRouter()
	router.Use(mw.LoggingMiddleware)

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", portNo),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
