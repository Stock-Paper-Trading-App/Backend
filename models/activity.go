package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Activity struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Symbol       string             `json:"symbol,omitempty" bson:"symbol,omitempty"`
	CompanyName  string             `json:"companyName,omitempty" bson:"companyName,omitempty"`
	Quantity     int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	Side         string             `json:"side,omitempty" bson:"side,omitempty"`
	Price        float64            `json:"price,omitempty" bson:"price,omitempty"`
	Initiated_on primitive.DateTime `json:"initiated_on,omitempty" bson:"initiated_on,omitempty"`
	User_id      primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

// include totalValue = quant*price in front end
