package handlers

import (
	"fictional-public-library/contracts"
	"fictional-public-library/errors"
	"fictional-public-library/literals"
	"fictional-public-library/services"
	"fictional-public-library/tracing"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

func DeleteBookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		var response *contracts.DeleteBookResponse
		var requestData *contracts.DeleteBookRequest

		ctx := request.Context()
		log := tracing.GetTracedLogEntry(ctx)

		err := UnmarshalRequest(request, &requestData, log)
		if err != nil {
			response = &contracts.DeleteBookResponse{
				Status: errors.FailureStatus,
				Errors: []*errors.ResponseError{
					&errors.UnmarshalRequestError,
				},
			}
			WriteToResponse(ctx, log, w, response, http.StatusBadRequest)
			return
		}

		log.WithFields(logrus.Fields{
			literals.RequestData: fmt.Sprintf("%+v", requestData),
		}).Info("request data")

		service := services.GetDeleteBookService()

		validationErrors := service.ValidateRequest(requestData)
		if len(validationErrors) > 0 {
			response = &contracts.DeleteBookResponse{
				Status: errors.FailureStatus,
				Errors: validationErrors,
			}
			WriteToResponse(ctx, log, w, response, http.StatusBadRequest)
			return
		}

		response, processError := service.ProcessRequest(ctx, requestData)
		if processError != nil {
			log.WithFields(logrus.Fields{
				literals.LLErrorMsg: processError.Message,
			}).Error("Error processing request")
			response = &contracts.DeleteBookResponse{
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
