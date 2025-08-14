package contracts

import "fictional-public-library/errors"

// AddBookResponse response body for POST /public-library/book/add API
type AddBookResponse struct {
	BookId string                  `contracts:"book_id,omitempty"`
	Status errors.ResponseStatus   `contracts:"status"`
	Errors []*errors.ResponseError `contracts:"errors,omitempty"`
}
