package dao

import (
	"context"
	"fictional-public-library/errors"
	"fictional-public-library/literals"
	"fictional-public-library/models"
	"fictional-public-library/routerconfig"
	"fictional-public-library/tracing"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"
)

var libraryDAO LibraryDAO
var libraryDAOOnce sync.Once

type libraryDAOStruct struct {
	config *routerconfig.RouterConfig
}

func InitLibraryDAO(conf *routerconfig.RouterConfig) LibraryDAO {
	libraryDAOOnce.Do(func() {
		libraryDAO = &libraryDAOStruct{config: conf}
	})
	return libraryDAO
}

func GetLibraryDAO() LibraryDAO {
	if libraryDAO == nil {
		panic("LibraryDAO not initialized")
	}
	return libraryDAO
}

type LibraryDAO interface {
	AddBook(ctx context.Context, book models.Book) (string, *errors.ResponseError)
	DeleteBook(ctx context.Context, bookID string) *errors.ResponseError
	UpdateBook(ctx context.Context, book models.Book) (string, *errors.ResponseError)
	RentBook(ctx context.Context, bookID string, rentedBy string) *errors.ResponseError
	FetchAllBooks(ctx context.Context) ([]*models.Book, *errors.ResponseError)
}

func (l libraryDAOStruct) AddBook(ctx context.Context, book models.Book) (string, *errors.ResponseError) {
	log := tracing.GetTracedLogEntry(ctx)

	book.Id = uuid.New().String()

	_, err := l.config.Mongo.Collection(literals.BookCollectionName).InsertOne(ctx, book)
	if err != nil {
		log.Error(err.Error(), "Error inserting book")
		return literals.EmptyString, &errors.AddBookError
	}

	return book.Id, nil
}

func (l libraryDAOStruct) DeleteBook(ctx context.Context, bookID string) *errors.ResponseError {
	log := tracing.GetTracedLogEntry(ctx)

	filter := bson.M{"_id": bookID}

	_, err := l.config.Mongo.Collection(literals.BookCollectionName).DeleteOne(ctx, filter)
	if err != nil {
		log.Error(err.Error(), "Error deleting book")
		return &errors.DeleteBookError
	}

	return nil
}

func (l libraryDAOStruct) UpdateBook(ctx context.Context, book models.Book) (string, *errors.ResponseError) {
	log := tracing.GetTracedLogEntry(ctx)

	filter := bson.M{"_id": book.Id}

	var setData = bson.M{
		"title":           book.Title,
		"author":          book.Author,
		"description":     book.Description,
		"publicationName": book.PublicationName,
		"updatedBy":       book.UpdatedBy,
		"updatedAt":       time.Now().UTC(),
	}
	var update = bson.M{
		"$set": setData,
	}

	_, err := l.config.Mongo.Collection(literals.BookCollectionName).UpdateOne(ctx, filter, update)
	if err != nil {
		log.Error(err.Error(), "Error updating book")
		return literals.EmptyString, &errors.UpdateBookError
	}

	return book.Id, nil
}

func (l libraryDAOStruct) RentBook(ctx context.Context, bookID string, rentedBy string) *errors.ResponseError {
	log := tracing.GetTracedLogEntry(ctx)

	filter := bson.M{"_id": bookID}
	var setData = bson.M{
		"rentedBy": rentedBy,
		"rentedAt": time.Now().UTC(),
	}
	var update = bson.M{
		"$set": setData,
	}

	_, err := l.config.Mongo.Collection(literals.BookCollectionName).FindOneAndUpdate(ctx, filter, update)
	if err != nil {
		log.Error(err.Error(), "Error renting book")
		return &errors.RentBookError
	}

	return nil
}

func (l libraryDAOStruct) FetchAllBooks(ctx context.Context) ([]*models.Book, *errors.ResponseError) {
	log := tracing.GetTracedLogEntry(ctx)
	filter := bson.M{}

	resultCursor, err := l.config.Mongo.Collection(literals.BookCollectionName).Find(ctx, filter)
	if err != nil {
		log.Error(err.Error(), "Error in fetching all books")
		return nil, &errors.FetchAllBooksError
	}

	defer resultCursor.Close(ctx)
	var books []*models.Book
	for resultCursor.Next(ctx) {
		var book models.Book
		err := resultCursor.Decode(&book)
		if err != nil {
			log.Error(err.Error(), "Error in decoding book record")
			return nil, &errors.BookDecodeError
		}
		books = append(books, &book)
	}
	
	return books, nil
}
