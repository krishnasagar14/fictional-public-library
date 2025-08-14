package models

import (
	"time"
)

var (
	BookSeqNo     int = 0
	RentBookSeqNo int = 0
)

type Book struct {
	Id              string    `json:"id" bson:"_id"`
	Title           string    `json:"title" bson:"title"`
	Author          string    `json:"author" bson:"author"`
	Description     string    `json:"description" bson:"description"`
	PublicationName string    `json:"publicationName" bson:"publicationName"`
	CreatedBy       string    `json:"createdBy" bson:"createdBy"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedBy       string    `json:"updatedBy" bson:"updatedBy"`
	UpdatedAt       time.Time `json:"updatedAt" bson:"updatedAt"`
	RentedBy        string    `json:"rentedBy" bson:"rentedBy"`
	RentedAt        time.Time `json:"rentedAt" bson:"rentedAt"`
}
