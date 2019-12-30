package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Car structure
type Car struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Model string             `json:"model,omitempty" bson:"model,omitempty"`
	Price string             `json:"price,omitempty" bson:"price,omitempty"`
}
