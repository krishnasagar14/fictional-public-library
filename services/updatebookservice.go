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
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type updateBookService struct {
	routerConfig *routerconfig.RouterConfig
	libraryDAO   dao.LibraryDAO
}

var updateBookSvcStruct UpdateBookServiceInterface
var updateBookServiceOnce sync.Once

// InitUpdateBookService ...
func InitUpdateBookService(config *routerconfig.RouterConfig, libDAO dao.LibraryDAO) UpdateBookServiceInterface {
	updateBookServiceOnce.Do(func() {
		updateBookSvcStruct = &updateBookService{
			routerConfig: config,
			libraryDAO:   libDAO,
		}
	})
	return updateBookSvcStruct
}

// GetUpdateBookService ...
func GetUpdateBookService() UpdateBookServiceInterface {
	if updateBookSvcStruct == nil {
		panic("updateBookService not initialized")
	}
	return updateBookSvcStruct
}

type UpdateBookServiceInterface interface {
	ValidateRequest(req *contracts.UpdateBookRequest) (validationErrors []*errors.ResponseError)
	ProcessRequest(ctx context.Context, req *contracts.UpdateBookRequest) (*contracts.UpdateBookResponse, *errors.ResponseError)
}

func (u updateBookService) ValidateRequest(req *contracts.UpdateBookRequest) (validationErrors []*errors.ResponseError) {
	if err := uuid.Validate(req.BookID); err != nil {
		validationErrors = append(validationErrors, &errors.InValidBookID)
	}

	if req.Title == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyTitleError)
	}

	if req.Description == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyDescriptionError)
	}

	if req.Author == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyAuthorError)
	}

	if req.UpdatedBy == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyCreatedByError)
	}

	return
}

func (u updateBookService) ProcessRequest(ctx context.Context, req *contracts.UpdateBookRequest) (*contracts.UpdateBookResponse, *errors.ResponseError) {
	log := tracing.GetTracedLogEntry(ctx)

	book := models.Book{
		Id:              req.BookID,
		Title:           req.Title,
		Description:     req.Description,
		Author:          req.Author,
		PublicationName: req.PublicationName,
		UpdatedBy:       req.UpdatedBy,
		UpdatedAt:       time.Now().UTC(),
	}

	_, err := u.libraryDAO.UpdateBook(ctx, book)
	if err != nil {
		log.WithFields(logrus.Fields{
			literals.LLErrorMsg: err.Error(),
		}).Error("Error updating book")
		return nil, err
	}

	resp := &contracts.UpdateBookResponse{
		BookId: book.Id,
		Status: errors.SuccessStatus,
	}

	return resp, nil
}
