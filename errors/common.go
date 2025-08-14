package errors

type ResponseError struct {
	error
	Message string `contracts:"message,omitempty"`
	Code    string `contracts:"code,omitempty"`
}

type ResponseStatus string

var (
	SuccessStatus ResponseStatus = "success"
	FailureStatus ResponseStatus = "failure"
)
