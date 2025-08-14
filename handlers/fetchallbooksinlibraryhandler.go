package handlers

import (
	"fictional-public-library/contracts"
	"fictional-public-library/errors"
	"fictional-public-library/literals"
	"fictional-public-library/services"
	"fictional-public-library/tracing"
	"github.com/sirupsen/logrus"
	"net/http"
)

func FetchAllBooksInLibraryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var response *contracts.FetchAllBooksInLibraryResponse

		ctx := request.Context()
		log := tracing.GetTracedLogEntry(ctx)

		log.Info("fetch all books in library called")

		service := services.GetFetchAllBooksBookService()

		response, processError := service.ProcessRequest(ctx)
		if processError != nil {
			log.WithFields(logrus.Fields{
				literals.LLErrorMsg: processError.Message,
			}).Error("Error processing request")
			response = &contracts.FetchAllBooksInLibraryResponse{
				Status: errors.FailureStatus,
				Errors: []*errors.ResponseError{
					&errors.AddBookServiceError,
				},
			}
			WriteToResponse(ctx, log, w, response, http.StatusInternalServerError)
			return
		}

		WriteToResponse(ctx, log, w, response, http.StatusOK)
	}
}
