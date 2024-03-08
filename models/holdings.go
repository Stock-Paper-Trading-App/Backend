package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Holdings struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Symbol      string             `json:"symbol,omitempty" bson:"symbol,omitempty"`
	CompanyName string             `json:"companyName,omitempty" bson:"companyName,omitempty"`
	Quantity    int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	User_id     primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}
