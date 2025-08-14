package errors

const Code = "LIB"

var (
	UnmarshalRequestError = ResponseError{
		Message: "Failed to unmarshal request",
		Code:    Code + "_001",
	}

	EmptyTitleError = ResponseError{
		Message: "Empty title",
		Code:    Code + "_002",
	}

	EmptyDescriptionError = ResponseError{
		Message: "Empty description",
		Code:    Code + "_003",
	}

	EmptyAuthorError = ResponseError{
		Message: "Empty author",
		Code:    Code + "_004",
	}

	EmptyCreatedByError = ResponseError{
		Message: "Empty createdBy",
		Code:    Code + "_005",
	}

	AddBookServiceError = ResponseError{
		Message: "Failed to add book into library",
		Code:    Code + "_006",
	}

	InValidBookID = ResponseError{
		Message: "Invalid book ID",
		Code:    Code + "_007",
	}

	EmptyDeletedBy = ResponseError{
		Message: "Empty deletedBy found",
		Code:    Code + "_008",
	}

	EmptyPublicationName = ResponseError{
		Message: "Empty publication name",
		Code:    Code + "_009",
	}
)
