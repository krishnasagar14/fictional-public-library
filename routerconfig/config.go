package routerconfig

import (
	"fictional-public-library/config"
	"fictional-public-library/db"
)

// RouterConfig Struct which stores all app dependencies required by different components
type RouterConfig struct {
	WebServerConfig *config.WebServerConfig
	Mongo           db.DatabaseHelper
}
