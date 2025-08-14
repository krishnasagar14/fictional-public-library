package handlers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fictional-public-library/literals"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	base64PaddingChar = "="
	secondsInOneYear  = "31536000"
)

func base64Decode(data []byte, log *logrus.Entry) ([]byte, error) {

	var b64StdDecError error
	var b64URLDecError error

	// Strip all padding characters ('=') from byte data, if any.
	// Found that any base64 encoded data can have at most 2 padding characters.
	data = bytes.TrimSuffix(data, []byte(base64PaddingChar))
	data = bytes.TrimSuffix(data, []byte(base64PaddingChar))

	// base64 decode using RawStdEncoding
	b64DecodedData := make([]byte, base64.RawStdEncoding.DecodedLen(len(data)))
	_, b64StdDecError = base64.RawStdEncoding.Decode(b64DecodedData, data)
	if b64StdDecError == nil {
		return b64DecodedData, nil
	}

	// base64 decode using RawURLEncoding
	b64DecodedData = make([]byte, base64.RawURLEncoding.DecodedLen(len(data)))
	_, b64URLDecError = base64.RawURLEncoding.Decode(b64DecodedData, data)
	if b64URLDecError == nil {
		return b64DecodedData, nil
	}

	log.Error("Base64 decode failed with both standard and url encoding techniques")
	log.Errorf("With RawStdEncoding, got error: %v", b64StdDecError.Error())
	log.Errorf("With RawURLEncoding, got error: %v", b64URLDecError.Error())
	return nil, errors.New("base 64 decode failed")
}

func UnmarshalRequest(r *http.Request, target interface{}, log *logrus.Entry) (err error) {
	var data []byte

	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		urlData, ok := r.URL.Query()[literals.RequestBody]
		if !ok || len(urlData[0]) < 1 {
			return nil
		}

		encodedReqData := []byte(urlData[0])
		var err error
		data, err = base64Decode(encodedReqData, log)
		if err != nil {
			return err
		}
	} else {
		ctx := r.Context()
		reqData := ctx.Value(literals.RequestBody)
		if reqData == nil {
			log.Error("Error while reading from request body: ", err)
			return errors.New("request body is nil")
		}

		data = reqData.([]byte)
	}

	err = json.Unmarshal(data, target)
	if err != nil {
		log.Error("Error while unmarshaling contracts request body: ", err)
		return err
	}

	return nil
}

func LogAndSendHTTPError(log *logrus.Entry, w http.ResponseWriter, error error,
	statusCode int) {
	log.WithError(error)
	http.Error(w, error.Error(), statusCode)
}

func SetResponseHeader(w http.ResponseWriter) {
	w.Header().Set(literals.HeaderContentType, "application/contracts")
	w.Header().Set("Content-Security-Policy", "default-src 'self';")
	w.Header().Set("Strict-Transport-Security", "max-age="+secondsInOneYear)
}

func IsSuccess(statusCode int) bool {
	return statusCode == 200 || statusCode == 204
}

func WriteToResponse(ctx context.Context, log *logrus.Entry, writer http.ResponseWriter, responseStruct interface{}, statusCode int) {
	var errResponse error
	if responseStruct != nil {
		response, err := json.Marshal(responseStruct)
		if err != nil {
			statusCode = http.StatusInternalServerError
			errResponse = errors.New("contracts marshal failed")

			LogAndSendHTTPError(log, writer, errResponse, http.StatusInternalServerError)
		}

		SetResponseHeader(writer)
		writer.WriteHeader(statusCode)
		_, err = writer.Write(response)
		if err != nil {
			statusCode = http.StatusInternalServerError
			errResponse = errors.New("write response failed")

			LogAndSendHTTPError(log, writer, errResponse, http.StatusInternalServerError)
		}
	}

	if IsSuccess(statusCode) {
		// Log the response code for the API
		log.WithFields(logrus.Fields{
			literals.LLEntryPoint:       ctx.Value(literals.EntrypointKey),
			literals.LLHTTPResponseCode: statusCode,
		}).Info(literals.APISucceeded)
	} else {
		log.WithFields(logrus.Fields{
			literals.LLEntryPoint:       ctx.Value(literals.EntrypointKey),
			literals.LLHTTPResponseCode: statusCode,
		}).Error(literals.APIFailed)
	}

}
