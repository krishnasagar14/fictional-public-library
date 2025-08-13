package db

import (
	"context"
	"errors"
	"fictional-public-library/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoErrorMap = map[int]error{
		11000: errors.New("duplicate key error"),
		121:   errors.New("document validation failed"),
	}
)

func (mc *mongoCollection) Count(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {

	count, err := mc.coll.CountDocuments(ctx, filter, opts...)

	if err == nil {
		return count, nil
	}

	logging.Log.Error("failed to get count, reason: ", err.Error())
	return -1, err
}

// FindOne executes a find command and returns a SingleResult for one document in the collection.
func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (SingleResultHelper, error) {

	singleResult := mc.coll.FindOne(ctx, filter, opts...)
	if err := singleResult.Err(); err != nil {
		// FindOne succeeded but returned no matching document
		if errors.Is(err, mongo.ErrNoDocuments) {
			logging.Log.Error("no document in result")
			return nil, mongo.ErrNoDocuments
		}
		logging.Log.Error("failed to find document in collection, reason: ", err.Error())
		return nil, err
	}
	return &mongoSingleResult{sr: singleResult}, nil
}

// Find executes a find command and returns a Cursor over the matching documents in the collection.
func (mc *mongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (CursorHelper, error) {

	cur, err := mc.coll.Find(ctx, filter, opts...)
	if err != nil {
		logging.Log.Error("failed to find documents, reason: ", err.Error())
		return nil, err
	}
	return &mongoCursor{cur: cur}, nil
}

// InsertOne executes an insert command to insert a single document into the collection.
func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {

	id, err := mc.coll.InsertOne(ctx, document, opts...)

	var writeException mongo.WriteException
	if errors.As(err, &writeException) {
		for _, mErr := range writeException.WriteErrors {
			if mongoError, ok := mongoErrorMap[mErr.Code]; ok {
				logging.Log.Error("failed to insert record in Mongo DB, reason: ", err.Error())
				return nil, mongoError
			}
		}
	}

	if err != nil {
		logging.Log.Error("failed to insert record in Mongo DB, reason: ", err.Error())
		return nil, err
	}
	return id, nil
}

// UpdateOne executes an update command to update at most one document in the collection.
func (mc *mongoCollection) UpdateOne(ctx context.Context, filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {

	res, err := mc.coll.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		logging.Log.Error("failed to update document in collection, reason: ", err.Error())
		return nil, err
	}
	return res, nil
}

// FindOneAndUpdate find and update 1 document also returns old document
func (mc *mongoCollection) FindOneAndUpdate(ctx context.Context, filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*mongo.SingleResult, error) {

	res := mc.coll.FindOneAndUpdate(ctx, filter, update, opts...)
	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			logging.Log.Error("no document in result")
			return nil, res.Err()
		}
		logging.Log.Error("failed to FindOneAndUpdate document in collection, reason: ", res.Err())
		return nil, res.Err()
	}
	return res, nil
}

// Decode will unmarshal the document represented by this SingleResult into v.
func (sr *mongoSingleResult) Decode(v interface{}) error {
	err := sr.sr.Decode(v)
	if err != nil {
		logging.Log.Error("failed to decode single result, reason: ", err.Error())
		return err
	}
	return nil
}

// Err returns the error from the operation that created this SingleResult.
func (sr *mongoSingleResult) Err() error {
	return sr.sr.Err()
}

// Next gets the next document for this cursor.
// It returns true if there were no errors and the cursor has not been exhausted.
func (c *mongoCursor) Next(ctx context.Context) bool {
	return c.cur.Next(ctx)
}

// Decode will unmarshal the current document into val and return any errors from the unmarshalling
// process without any modification.
func (c *mongoCursor) Decode(val interface{}) error {
	err := c.cur.Decode(val)
	if err != nil {
		logging.Log.Error("failed to decode, reason: ", err.Error())
		return err
	}
	return nil
}

// Close closes this cursor.
func (c *mongoCursor) Close(ctx context.Context) error {
	err := c.cur.Close(ctx)
	if err != nil {
		logging.Log.Error("closing cursor failed, reason: ", err.Error())
		return err
	}
	return nil
}

// Err returns the last error seen by the Cursor, or nil if no error has occurred.
func (c *mongoCursor) Err() error {
	err := c.cur.Err()
	if err != nil {
		logging.Log.Error("last error seen by cursor, reason: ", err.Error())
		return err
	}
	return nil
}
