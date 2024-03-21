package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Networth struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Networth     float64            `json:"networth,omitempty" bson:"networth,omitempty"`
	Initiated_on primitive.DateTime `json:"initiated_on,omitempty" bson:"initiated_on,omitempty"`
	User_id      primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}
