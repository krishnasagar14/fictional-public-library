package server

import (
	"fictional-public-library/config"
	"fictional-public-library/dao"
	"fictional-public-library/db"
	"fictional-public-library/literals"
	"fictional-public-library/logging"
	"fictional-public-library/routerconfig"
	"fictional-public-library/services"
	"fmt"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
)

type Server struct {
	Configuration *config.WebServerConfig
	Router        *Router
}

// NewServer ...
func NewServer(config *config.WebServerConfig) *Server {
	return &Server{
		Configuration: config,
		Router:        NewRouter(),
	}
}

func RunServer() (err error) {
	webServerConfig, logCfg, err := config.FromEnv()
	if err != nil {
		return err
	}

	err = logging.Initialize(logCfg)
	if err != nil {
		return err
	}

	routerCfg := &routerconfig.RouterConfig{
		WebServerConfig: webServerConfig,
	}
	dao.InitLibraryDAO(routerCfg)
	initServices(routerCfg)

	// initialize mongo
	if err = mongoInit(webServerConfig, routerCfg); err != nil {
		return err
	}

	server := NewServer(webServerConfig)
	server.Router.InitializeRouter(routerCfg)

	crossOriginServer := corsSetup()

	fmt.Printf("Starting HTTP server on port %s", webServerConfig.Port)
	err = http.ListenAndServe(":"+webServerConfig.Port, crossOriginServer.Handler(server.Router))
	if err != nil {
		return err
	}

	return nil
}

func initServices(config *routerconfig.RouterConfig) {
	services.InitAddBookService(config, dao.GetLibraryDAO())
}

func mongoInit(wc *config.WebServerConfig, rc *routerconfig.RouterConfig) error {
	// client & database options
	clientOption := &options.ClientOptions{
		MaxPoolSize: &wc.MongoMaxPoolSize,
		MinPoolSize: &wc.MongoMinPoolSize,
	}

	clOpts := []*options.ClientOptions{clientOption}
	var dbOpts []*options.DatabaseOptions

	// connect to database
	mongoConnManager, err := db.NewDatabase(wc.MongoCfg, clOpts, dbOpts)
	if err != nil {
		logging.Log.Errorf("MongoDB initialization failed, reason: %v", err.Error())
		return err
	}

	// set db handle in router config struct
	rc.Mongo = mongoConnManager
	return nil
}

func corsSetup() *cors.Cors {

	var allowedOrigins []string
	allowedOrigins = append(allowedOrigins, "*")

	// TODO: add authorization header

	return cors.New(cors.Options{
		Logger:         logging.Log,
		Debug:          true,
		AllowedOrigins: allowedOrigins,
		AllowedHeaders: []string{literals.HeaderAccept,
			literals.HeaderContentType,
			literals.HeaderContentLength,
			literals.HeaderAcceptEncoding,
			literals.HeaderAccessControlAllowOrigin,
		},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodPost, http.MethodGet, http.MethodOptions, http.MethodPut, http.MethodDelete, http.MethodPatch},
	})
}
