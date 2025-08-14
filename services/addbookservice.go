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
	"time"
)

type addBookService struct {
	routerConfig *routerconfig.RouterConfig
	libraryDAO   dao.LibraryDAO
}

type AddBookServiceInterface interface {
	ValidateRequest(req *contracts.AddBookRequest) (validationErrors []*errors.ResponseError)
	ProcessRequest(ctx context.Context, req *contracts.AddBookRequest) (*contracts.AddBookResponse, *errors.ResponseError)
}

var addBookSvcStruct AddBookServiceInterface
var addBookServiceOnce sync.Once

// InitAddBookService ...
func InitAddBookService(config *routerconfig.RouterConfig, libDAO dao.LibraryDAO) AddBookServiceInterface {
	addBookServiceOnce.Do(func() {
		addBookSvcStruct = &addBookService{
			routerConfig: config,
			libraryDAO:   libDAO,
		}
	})
	return addBookSvcStruct
}

// GetAddBookService ...
func GetAddBookService() AddBookServiceInterface {
	if addBookSvcStruct == nil {
		panic("AddBookService not initialized")
	}
	return addBookSvcStruct
}

func (a addBookService) ValidateRequest(req *contracts.AddBookRequest) (validationErrors []*errors.ResponseError) {
	if req.Title == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyTitleError)
	}

	if req.Description == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyDescriptionError)
	}

	if req.Author == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyAuthorError)
	}

	if req.CreatedBy == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyCreatedByError)
	}

	if req.PublicationName == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyPublicationName)
	}

	return
}

func (a addBookService) ProcessRequest(ctx context.Context, req *contracts.AddBookRequest) (*contracts.AddBookResponse, *errors.ResponseError) {

	log := tracing.GetTracedLogEntry(ctx)

	book := models.Book{
		Title:           req.Title,
		Description:     req.Description,
		Author:          req.Author,
		CreatedBy:       req.CreatedBy,
		PublicationName: req.PublicationName,
		CreatedAt:       time.Now().UTC(),
	}

	bookID, err := a.libraryDAO.AddBook(ctx, book)
	if err != nil {
		log.WithFields(logrus.Fields{
			literals.LLErrorMsg: err.Error(),
		}).Error("Error adding book")
		return nil, err
	}

	resp := &contracts.AddBookResponse{
		Status: errors.SuccessStatus,
		BookId: bookID,
	}

	return resp, nil
}
