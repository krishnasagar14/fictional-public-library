package models

import (
	"strconv"
	"time"
)

var (
	BookSeqNo     int = 0
	RentBookSeqNo int = 0
)

type Book struct {
	Id          string    `contracts:"id" bson:"_id"`
	Title       string    `contracts:"title" bson:"title"`
	Author      string    `contracts:"author" bson:"author"`
	Description string    `contracts:"description" bson:"description"`
	CreatedBy   string    `contracts:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time `contracts:"createdAt" bson:"createdAt"`
}

func GetBookSeqNumber() string {
	BookSeqNo++
	return strconv.Itoa(BookSeqNo)
}

type RentBook struct {
	Id        string    `contracts:"id" bson:"_id"`
	BookID    string    `contracts:"book_id" bson:"book_id"`
	CreatedOn time.Time `contracts:"createdOn" bson:"createdOn"`
	ExpiresOn time.Time `contracts:"expiresOn" bson:"expiresOn"`
	UserName  string    `contracts:"user_name" bson:"user_name"`
	RentCost  float64   `contracts:"rent_cost" bson:"rent_cost"`
	Status    string    `contracts:"status" bson:"status"`
}

func GetRentBookSeqNumber() string {
	RentBookSeqNo++
	return strconv.Itoa(RentBookSeqNo)
}
