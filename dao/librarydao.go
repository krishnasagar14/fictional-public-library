package dao

import (
	"context"
	"fictional-public-library/errors"
	"fictional-public-library/literals"
	"fictional-public-library/models"
	"fictional-public-library/routerconfig"
	"fictional-public-library/tracing"
	"github.com/google/uuid"
	"sync"
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
