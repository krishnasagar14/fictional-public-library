package errors

const code = "DB"

var (
	AddBookError = ResponseError{
		Code:    code + "_001",
		Message: "add book error",
	}

	DeleteBookError = ResponseError{
		Code:    code + "_002",
		Message: "delete book error",
	}

	UpdateBookError = ResponseError{
		Code:    code + "_003",
		Message: "update book error",
	}

	RentBookError = ResponseError{
		Code:    code + "_004",
		Message: "rent book error",
	}

	FetchAllBooksError = ResponseError{
		Code:    code + "_005",
		Message: "fetch all books error",
	}

	BookDecodeError = ResponseError{
		Code:    code + "_006",
		Message: "book record decode error",
	}
)
