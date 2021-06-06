package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID       primitive.ObjectID  `bson:"_id" json:"id,omitempty"`
    Name     string `json:"name"`
    Location string `json:"location"`
    Age      int64  `json:"age"`
}

type Response struct {
    ID      string  `json:"id,omitempty"`
    Message string `json:"message,omitempty"`
}
