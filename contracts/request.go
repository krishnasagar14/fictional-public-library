package contracts

// AddBookRequest request body for POST /public-library/v1/book/add API
type AddBookRequest struct {
	Title           string `json:"title"`           // Required; title of book
	Author          string `json:"author"`          // Required; author name of book
	Description     string `json:"description"`     // Required; Book brief description
	PublicationName string `json:"publicationName"` // required; publication name of book
	CreatedBy       string `json:"createdBy" `      // Required; username is expected value
}

// DeleteBookRequest request body for POST /public-library/v1/book/delete API
type DeleteBookRequest struct {
	BookID    string `json:"bookID"`    // Required; uuid 4 value of book ID
	DeletedBy string `json:"deletedBy"` // Required; username who deleted the book
}

// UpdateBookRequest request body for PUT /public-library/v1/book/update API
type UpdateBookRequest struct {
	BookID          string `json:"bookID"`          // Required; uuid 4 value of book ID
	Title           string `json:"title"`           // Required; title of book
	Author          string `json:"author"`          // Required; author name of book
	Description     string `json:"description"`     // Required; Book brief description
	PublicationName string `json:"publicationName"` // required; publication name of book
	UpdatedBy       string `json:"updatedBy" `      // Required; username is expected value
}

// RentBookRequest request body for PATCH /public-library/book/rent API
type RentBookRequest struct {
	BookID   string `json:"bookID,omitempty"`   // Required; uuid 4 value of book ID
	RentedBy string `json:"rentedBy,omitempty"` // Required; username is expected value
}
