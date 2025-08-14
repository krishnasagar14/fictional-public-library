package contracts

import (
	"fictional-public-library/errors"
	"time"
)

// AddBookResponse response body for POST /public-library/book/add API
type AddBookResponse struct {
	BookId string                  `json:"bookId,omitempty"`
	Status errors.ResponseStatus   `json:"status"`
	Errors []*errors.ResponseError `json:"errors,omitempty"`
}

// DeleteBookResponse response body for POST /public-library/v1/book/delete API
type DeleteBookResponse struct {
	Status errors.ResponseStatus   `json:"status"`
	Errors []*errors.ResponseError `json:"errors,omitempty"`
}

// UpdateBookResponse response body for PUT /public-library/book/update API
type UpdateBookResponse struct {
	BookId string                  `json:"bookId,omitempty"`
	Status errors.ResponseStatus   `json:"status"`
	Errors []*errors.ResponseError `json:"errors,omitempty"`
}

// RentBookResponse response body for PATCH /public-library/book/rent API
type RentBookResponse struct {
	BookId string                  `json:"bookId,omitempty"`
	Status errors.ResponseStatus   `json:"status"`
	Errors []*errors.ResponseError `json:"errors,omitempty"`
}

type BookResponse struct {
	BookId          string    `json:"bookId"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Description     string    `json:"description"`
	PublicationName string    `json:"publicationName"`
	CreatedBy       string    `json:"createdBy"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedBy       string    `json:"updatedBy"`
	UpdatedAt       time.Time `json:"updatedAt"`
	RentedBy        string    `json:"rentedBy"`
	RentedAt        time.Time `json:"rentedAt"`
}

// FetchAllBooksInLibraryResponse response body for GET /public-library/books API
type FetchAllBooksInLibraryResponse struct {
	Books  []*BookResponse         `json:"books,omitempty"`
	Status errors.ResponseStatus   `json:"status"`
	Errors []*errors.ResponseError `json:"errors,omitempty"`
}
