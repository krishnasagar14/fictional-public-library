package contracts

// AddBookRequest request body for POST /public-library/book/add API
type AddBookRequest struct {
	Title       string `contracts:"title"`                          // Required
	Author      string `contracts:"author"`                         // Required
	Description string `contracts:"description" bson:"description"` // Required
	CreatedBy   string `contracts:"createdBy" bson:"createdBy"`     // Required
}
