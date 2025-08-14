package config

import (
	"fictional-public-library/logging"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
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

// formConnectionURI builds mongo connection URI.
// Reference doc - https://docs.mongodb.com/manual/reference/connection-string/ not limited to
func (cnf *MongoConfig) formConnectionURI() string {

	var authMode bool // when username & password is present, auth must be done
	if cnf.User != "" && cnf.Password != "" {
		authMode = true
	}

	var mURI string
	mURI = cnf.URIPrefix + "://"

	if authMode {
		mURI += cnf.User + ":" + url.QueryEscape(cnf.Password) + "@"
	}

	mURI += cnf.Host

	if authMode {
		if cnf.AuthSource != "" {
			mURI += "/" + cnf.AuthSource
		} else {
			logging.Log.Warn("Authsource should be specified when username and password are present")
		}
	}

	if len(cnf.Flags) != 0 {
		mURI += "/" // trailing slash must be present before query params
		u, _ := url.Parse(mURI)
		q := u.Query()
		for optionKey, optionVal := range cnf.Flags {
			q.Add(optionKey, optionVal)
		}
		u.RawQuery = q.Encode()
		mURI = u.String()
	}

	return mURI
}

func (cnf *MongoConfig) SetClientOptions(clOpts ...*options.ClientOptions) []*options.ClientOptions {
	if clOpts == nil {
		clOpts = []*options.ClientOptions{}
	}

	mURI := cnf.formConnectionURI()

	clOpt := options.Client().ApplyURI(mURI)
	clOpts = append(clOpts, clOpt)

	return clOpts
}

type WebServerConfig struct {
	Port             string       `required:"true" split_words:"true" default:"50051"`
	MongoCfg         *MongoConfig `required:"true" split_words:"true"`
	MongoMaxPoolSize uint64       `split_words:"true" default:"10"`
	MongoMinPoolSize uint64       `split_words:"true" default:"5"`
	RoutePrefix      string       `required:"false" split_words:"true" default:"/public-library"`
}

func FromEnv() (cfg *WebServerConfig, lcfg *logging.LogConfig, err error) {
	fromFileToEnv()

	cfg = &WebServerConfig{}

	err = envconfig.Process(envPrefix, cfg)
	if err != nil {
		return nil, nil, err
	}

	lcfg = &logging.LogConfig{}
	err = envconfig.Process("", lcfg)
	if err != nil {
		return nil, nil, err
	}

	return cfg, lcfg, nil
}

func fromFileToEnv() {

	cfgFilename := "./config/config.localhost.env"

	err := godotenv.Load(cfgFilename)
	if err != nil {
		fmt.Println("No config files found to load to env. Defaulting to environment.", err.Error())
	}

}
