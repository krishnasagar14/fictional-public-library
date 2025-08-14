package db

import (
	"context"
	"fictional-public-library/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseHelper interface {
	Collection(string, ...*options.CollectionOptions) CollectionHelper
	Client() ClientHelper
}

// mongoDatabase implements DatabaseHelper
type mongoDatabase struct {
	db  *mongo.Database
	cfg *config.MongoConfig
}

// mongoCollection implements CollectionHelper
type mongoCollection struct {
	coll *mongo.Collection
}

// mongoCursor implements CursorHelper
type mongoCursor struct {
	cur *mongo.Cursor
}

// mongoSingleResult implements SingleResultHelper
type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type SingleResultHelper interface {
	Decode(interface{}) error
	Err() error
}

type CursorHelper interface {
	Next(context.Context) bool
	Decode(interface{}) error
	Close(context.Context) error
	Err() error
}

// mongoClient implements ClientHelper
type mongoClient struct {
	cl                *mongo.Client
	connectTimeout    int
	sessionTimeout    int
	enableTransaction bool
}

// ClientHelper is the interface that wraps client methods.
type ClientHelper interface {
	Ping() error
	Connect(context.Context) error
}

// CollectionHelper is the interface that wraps collection methods.
type CollectionHelper interface {
	Count(context.Context, interface{}, ...*options.CountOptions) (int64, error)
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) (SingleResultHelper, error)
	Find(context.Context, interface{}, ...*options.FindOptions) (CursorHelper, error)
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	UpdateOne(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
	FindOneAndUpdate(ctx context.Context, filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error)
}
