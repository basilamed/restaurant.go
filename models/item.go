package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Price	    float64            `json:"price,omitempty" bson:"price,omitempty"`
	AvailableUnits int             `json:"availableUnits,omitempty" bson:"availableUnits,omitempty"`
}