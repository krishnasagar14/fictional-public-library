package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "LIBRARY"

type MongoConfig struct {
	Name                  string            `required:"true" split_words:"true"`
	AuthSource            string            `required:"false" split_words:"true"`
	URIPrefix             string            `required:"false" split_words:"true" default:"mongodb+srv"`
	User                  string            `split_words:"true"`
	Password              string            `split_words:"true"`
	Host                  string            `required:"true" split_words:"true"`
	Flags                 map[string]string `split_words:"true"`
	ConnectRetryCount     int               `split_words:"true" default:"3"`
	ConnectRetryInterval  int               `split_words:"true" default:"5"`
	ConnectTimeout        int               `split_words:"true" default:"15"`
	SessionTimeout        int               `split_words:"true" default:"15"`
	EnableTransaction     bool              `split_words:"true" default:"false"`
	QueryImportTimeoutSec int               `split_words:"true" default:"30"`
}

type WebServerConfig struct {
}

type LogConfig struct {
	FileEnabled    bool   `envconfig:"LOG_FILE_ENABLED" default:"false"`
	FileLocation   string `envconfig:"LOG_FILE_NAME" default:"server.log"`
	ConsoleEnabled bool   `envconfig:"LOG_CONSOLE_ENABLED" default:"true"`
	JSONFormat     bool   `envconfig:"LOG_JSON_FORMAT"`
}

func FromEnv() (cfg *WebServerConfig, lcfg *LogConfig, err error) {
	fromFileToEnv()

	cfg = &WebServerConfig{}

	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, nil, err
	}

	lcfg = &LogConfig{}
	err = envconfig.Process("", lcfg)
	if err != nil {
		return nil, nil, err
	}

	return cfg, lcfg, nil
}

func fromFileToEnv() {

	cfgFilename := "./config.localhost.env"

	err := godotenv.Load(cfgFilename)
	if err != nil {
		fmt.Println("No config files found to load to env. Defaulting to environment.")
	}

}
