package errors

const code = "DB"

var (
	AddBookError = ResponseError{
		Code:    code + "_001",
		Message: "add book error",
	}
)
