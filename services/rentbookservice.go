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

type rentBookService struct {
	routerConfig *routerconfig.RouterConfig
	libraryDAO   dao.LibraryDAO
}

var rentBookSvcStruct RentBookServiceInterface
var rentBookServiceOnce sync.Once

// InitRentBookService ...
func InitRentBookService(config *routerconfig.RouterConfig, libDAO dao.LibraryDAO) RentBookServiceInterface {
	rentBookServiceOnce.Do(func() {
		rentBookSvcStruct = &rentBookService{
			routerConfig: config,
			libraryDAO:   libDAO,
		}
	})
	return rentBookSvcStruct
}

// GetRentBookService ...
func GetRentBookService() RentBookServiceInterface {
	if rentBookSvcStruct == nil {
		panic("rentBookService not initialized")
	}
	return rentBookSvcStruct
}

type RentBookServiceInterface interface {
	ValidateRequest(req *contracts.RentBookRequest) (validationErrors []*errors.ResponseError)
	ProcessRequest(ctx context.Context, req *contracts.RentBookRequest) (*contracts.RentBookResponse, *errors.ResponseError)
}

func (r rentBookService) ValidateRequest(req *contracts.RentBookRequest) (validationErrors []*errors.ResponseError) {
	if err := uuid.Validate(req.BookID); err != nil {
		validationErrors = append(validationErrors, &errors.InValidBookID)
	}

	if req.RentedBy == literals.EmptyString {
		validationErrors = append(validationErrors, &errors.EmptyDeletedBy)
	}

	return
}

func (r rentBookService) ProcessRequest(ctx context.Context, req *contracts.RentBookRequest) (*contracts.RentBookResponse, *errors.ResponseError) {
	log := tracing.GetTracedLogEntry(ctx)

	err := r.libraryDAO.RentBook(ctx, req.BookID, req.RentedBy)
	if err != nil {
		log.WithFields(logrus.Fields{
			literals.LLErrorMsg: err.Error(),
		}).Error("Error deleting book")
		return nil, err
	}

	resp := &contracts.RentBookResponse{
		BookId: req.BookID,
		Status: errors.SuccessStatus,
	}

	return resp, nil
}
