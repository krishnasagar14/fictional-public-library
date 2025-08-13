package db

import (
	"context"
	"errors"
	"fictional-public-library/config"
	"fictional-public-library/logging"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection gets a handle for a collection with the given name configured with the given CollectionOptions.
func (md *mongoDatabase) Collection(colName string, opts ...*options.CollectionOptions) CollectionHelper {
	collection := md.db.Collection(colName, opts...)
	return &mongoCollection{
		coll: collection,
	}
}

// Client returns the Client the Database was created from.
func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{
		cl:                client,
		connectTimeout:    md.cfg.ConnectTimeout,
		sessionTimeout:    md.cfg.SessionTimeout,
		enableTransaction: md.cfg.EnableTransaction,
	}
}

// NewDatabase accepts a client and returns a database handle.
func NewDatabase(cnf *config.MongoConfig, clOpts []*options.ClientOptions, dbOpts []*options.DatabaseOptions) (DatabaseHelper, error) {
	var err error
	err = sanitizeAndValidateConfig(cnf)
	if err != nil {
		return nil, err
	}
	clOpts = cnf.SetClientOptions(clOpts...)

	isConnectionSuccessful := false
	var mcl *mongoClient
	for i := 0; i < cnf.ConnectRetryCount; i++ {

		logging.Log.Infof("Attempt #%d: Connecting to MongoDB...", i+1)
		mcl, err = attemptConnection(cnf, clOpts...)
		if err != nil {
			// failed last attempt, break loop
			if i == cnf.ConnectRetryCount-1 {
				break
			} else {
				logging.Log.Infof("Reconnection will be attempted in %d seconds", cnf.ConnectRetryInterval)
				time.Sleep(time.Duration(cnf.ConnectRetryInterval) * time.Second)
				continue
			}
		} else {
			isConnectionSuccessful = true
			break
		}
	}

	if !isConnectionSuccessful {
		return nil, err
	}

	return &mongoDatabase{
		db:  mcl.cl.Database(cnf.Name, dbOpts...),
		cfg: cnf,
	}, nil
}

func attemptConnection(cnf *config.MongoConfig, clOpts ...*options.ClientOptions) (*mongoClient, error) {
	// new mongo client
	c, err := mongo.NewClient(clOpts...)
	if err != nil {
		logging.Log.Error("failed to create new mongo client, reason: ", err.Error())
		return nil, errors.New("failed to create new mongo client, reason: " + err.Error())
	}

	mcl := &mongoClient{cl: c, connectTimeout: cnf.ConnectTimeout,
		sessionTimeout: cnf.SessionTimeout}

	connectCtx, cancel := context.WithTimeout(context.Background(),
		time.Duration(mcl.connectTimeout)*time.Second)
	defer cancel()

	err = mcl.Connect(connectCtx)
	if err != nil {
		return nil, err
	}

	return mcl, nil
}

func sanitizeAndValidateConfig(cfg *config.MongoConfig) error {
	// trim leading & trailing spaces for each config value
	cfg.Name = strings.TrimSpace(cfg.Name)
	cfg.URIPrefix = strings.TrimSpace(cfg.URIPrefix)
	cfg.User = strings.TrimSpace(cfg.User)
	cfg.Password = strings.TrimSpace(cfg.Password)
	cfg.Host = strings.TrimSpace(cfg.Host)
	cfg.Host = filepath.Clean(cfg.Host) // removes trailing slashes, if any
	for i := range cfg.Flags {
		cfg.Flags[i] = strings.TrimSpace(cfg.Flags[i])
	}

	// look for required config values that are missing
	if cfg.Name == "" {
		logging.Log.Error("Missing database name in configuration")
		return errors.New("missing database name in configuration")
	}
	if cfg.Host == "" {
		logging.Log.Error("Missing host url in configuration")
		return errors.New("missing host url in configuration")
	}

	// URI prefix should be `mongodb` or `mongodb+srv`
	if !(cfg.URIPrefix == "mongodb" || cfg.URIPrefix == "mongodb+srv") {
		logging.Log.Errorf("Invalid URI prefix %v was supplied", cfg.URIPrefix)
		return errors.New("invalid URI prefix " + cfg.URIPrefix + " was supplied")
	}

	if cfg.ConnectRetryCount <= 1 {
		logging.Log.Errorf("Retry count cannot be less than or equal to 1")
		return errors.New("retry count cannot be less than or equal to 1")
	}
	if cfg.ConnectRetryInterval <= 0 {
		logging.Log.Errorf("Retry interval should be greater than 0 seconds")
		return errors.New("retry interval should be greater than 0 seconds")
	}
	if cfg.ConnectTimeout <= 0 {
		logging.Log.Errorf("Connection timeout should be greater than 0 seconds")
		return errors.New("connection timeout should be greater than 0 seconds")
	}
	if cfg.SessionTimeout <= 0 {
		logging.Log.Errorf("Session timeout should be greater than 0 seconds")
		return errors.New("session timeout should be greater than 0 seconds")
	}

	return nil
}

// Connect establishes a connection to mongo instance.
func (mc *mongoClient) Connect(ctx context.Context) error {
	err := mc.cl.Connect(ctx)
	if err != nil {
		logging.Log.Error("failed to connect to Mongo DB, reason: ", err.Error())
		return err
	}

	if err = mc.Ping(); err != nil {
		return err
	}

	logging.Log.Info("Connected to Mongo DB!")
	return nil
}

// Ping sends a ping command to verify that the client can connect to the mongo deployment.
func (mc *mongoClient) Ping() error {
	var err error
	pingCtx, cancel := context.WithTimeout(context.Background(),
		time.Duration(mc.connectTimeout)*time.Second)
	defer cancel()

	err = mc.cl.Ping(pingCtx, nil)
	if err != nil {
		logging.Log.Error("failed to ping Mongo DB, reason: ", err.Error())
		return err
	}

	return nil
}
