package services

import (
	"context"
	"fictional-public-library/contracts"
	"fictional-public-library/dao"
	"fictional-public-library/errors"
	"fictional-public-library/literals"
	"fictional-public-library/models"
	"fictional-public-library/routerconfig"
	"fictional-public-library/tracing"
	"github.com/sirupsen/logrus"
	"sync"
)

type fetchAllBooksBookService struct {
	routerConfig *routerconfig.RouterConfig
	libraryDAO   dao.LibraryDAO
}

var fetchAllBooksBookSvcStruct FetchAllBooksBookServiceInterface
var fetchAllBooksBookServiceOnce sync.Once

// InitFetchAllBooksBookService ...
func InitFetchAllBooksBookService(config *routerconfig.RouterConfig, libDAO dao.LibraryDAO) FetchAllBooksBookServiceInterface {
	fetchAllBooksBookServiceOnce.Do(func() {
		fetchAllBooksBookSvcStruct = &fetchAllBooksBookService{
			routerConfig: config,
			libraryDAO:   libDAO,
		}
	})
	return fetchAllBooksBookSvcStruct
}

// GetFetchAllBooksBookService ...
func GetFetchAllBooksBookService() FetchAllBooksBookServiceInterface {
	if fetchAllBooksBookSvcStruct == nil {
		panic("fetchAllBooksBookService not initialized")
	}
	return fetchAllBooksBookSvcStruct
}

type FetchAllBooksBookServiceInterface interface {
	ProcessRequest(ctx context.Context) (*contracts.FetchAllBooksInLibraryResponse, *errors.ResponseError)
}

func (f fetchAllBooksBookService) ProcessRequest(ctx context.Context) (*contracts.FetchAllBooksInLibraryResponse, *errors.ResponseError) {
	log := tracing.GetTracedLogEntry(ctx)

	books, err := f.libraryDAO.FetchAllBooks(ctx)
	if err != nil {
		log.WithFields(logrus.Fields{
			literals.LLErrorMsg: err.Error(),
		}).Error("Error fetching all books")
		return nil, err
	}

	resp := &contracts.FetchAllBooksInLibraryResponse{
		Status: errors.SuccessStatus,
		Books:  prepareBookResponse(books),
	}

	return resp, nil
}

func prepareBookResponse(books []*models.Book) []*contracts.BookResponse {
	var bookResponses []*contracts.BookResponse
	for _, book := range books {
		bookResponse := &contracts.BookResponse{
			BookId:          book.Id,
			Title:           book.Title,
			Author:          book.Author,
			Description:     book.Description,
			PublicationName: book.PublicationName,
			RentedBy:        book.RentedBy,
			RentedAt:        book.RentedAt,
			UpdatedAt:       book.UpdatedAt,
			UpdatedBy:       book.UpdatedBy,
			CreatedAt:       book.CreatedAt,
			CreatedBy:       book.CreatedBy,
		}
		bookResponses = append(bookResponses, bookResponse)
	}
	return bookResponses
}
