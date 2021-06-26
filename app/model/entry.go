package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Buysell struct {
	Buy  string
	Sell string
}

type Rates struct {
	ER Buysell
	TT Buysell
	BN Buysell
}

type Exc_rates struct {
	Symbol string
	Rates  Rates
	// Rates string
	// Date time.Time
}

type Entry struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	Date       time.Time          `bson:"date" json:"date"`
	Exc_Rates  []Exc_rates        `bson:"exc_rates" json:"exc_rates"`
	Created_At time.Time          `bson:"created_at" json:"created_at"`
	Updated_At time.Time          `bson:"updated_at" json:"updated_at"`
}
