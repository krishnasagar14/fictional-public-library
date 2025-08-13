package db

import (
	"context"
	"fictional-public-library/logging"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

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

func ConnectDB(dbName string) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
		return err
	}
	// Double sure connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	db = client.Database(dbName)
	log.Println("Connected to database successfully !!!")
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
