package models

import (
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	BookSeqNo     int = 0
	RentBookSeqNo int = 0
)

type Book struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Author      string             `json:"author" bson:"author"`
	Description string             `json:"description" bson:"description"`
}

func GetBookSeqNumber() string {
	BookSeqNo++
	return strconv.Itoa(BookSeqNo)
}

type RentBook struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	BookID    string             `json:"book_id" bson:"book_id"`
	CreatedOn time.Time          `json:"createdOn" bson:"createdOn"`
	ExpiresOn time.Time          `json:"expiresOn" bson:"expiresOn"`
	UserName  string             `json:"user_name" bson:"user_name"`
	RentCost  float64            `json:"rent_cost" bson:"rent_cost"`
	Status    string             `json:"status" bson:"status"`
}

func GetRentBookSeqNumber() string {
	RentBookSeqNo++
	return strconv.Itoa(RentBookSeqNo)
}
