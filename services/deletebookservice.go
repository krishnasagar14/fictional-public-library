package services

import (
	"context"
	"fictional-public-library/contracts"
	"fictional-public-library/dao"
	"fictional-public-library/errors"
	"fictional-public-library/literals"
	"fictional-public-library/routerconfig"
	"fictional-public-library/tracing"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
)

type deleteBookService struct {
	routerConfig *routerconfig.RouterConfig
	libraryDAO   dao.LibraryDAO
}

var deleteBookSvcStruct DeleteBookServiceInterface
var deleteBookServiceOnce sync.Once

// InitDeleteBookService ...
func InitDeleteBookService(config *routerconfig.RouterConfig, libDAO dao.LibraryDAO) DeleteBookServiceInterface {
	deleteBookServiceOnce.Do(func() {
		deleteBookSvcStruct = &deleteBookService{
			routerConfig: config,
			libraryDAO:   libDAO,
		}
	})
	return deleteBookSvcStruct
}

// GetDeleteBookService ...
func GetDeleteBookService() DeleteBookServiceInterface {
	if deleteBookSvcStruct == nil {
		panic("deleteBookService not initialized")
	}
	return deleteBookSvcStruct
}

type DeleteBookServiceInterface interface {
	ValidateRequest(req *contracts.DeleteBookRequest) (validationErrors []*errors.ResponseError)
	ProcessRequest(ctx context.Context, req *contracts.DeleteBookRequest) (*contracts.DeleteBookResponse, *errors.ResponseError)
}

func (d deleteBookService) ValidateRequest(req *contracts.DeleteBookRequest) (validationErrors []*errors.ResponseError) {
	if err := uuid.Validate(req.BookID); err != nil {
		validationErrors = append(validationErrors, &errors.InValidBookID)
	}

	if req.DeletedBy == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyDeletedBy)
	}

	return
}

func (d deleteBookService) ProcessRequest(ctx context.Context, req *contracts.DeleteBookRequest) (*contracts.DeleteBookResponse, *errors.ResponseError) {
	log := tracing.GetTracedLogEntry(ctx)

	err := d.libraryDAO.DeleteBook(ctx, req.BookID)
	if err != nil {
		log.WithFields(logrus.Fields{
			literals.LLErrorMsg: err.Error(),
		}).Error("Error deleting book")
		return nil, err
	}

	resp := &contracts.DeleteBookResponse{
		Status: errors.SuccessStatus,
	}

	return resp, nil
}
