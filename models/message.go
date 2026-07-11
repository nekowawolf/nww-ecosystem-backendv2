package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Message struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Text1  string             `bson:"text_1" json:"text_1"`
	Text2  string             `bson:"text_2" json:"text_2"`
	Text3  string             `bson:"text_3" json:"text_3"`
	Text4  string             `bson:"text_4" json:"text_4"`
	Text5  string             `bson:"text_5" json:"text_5"`
	Text6  string             `bson:"text_6" json:"text_6"`
	Text7  string             `bson:"text_7" json:"text_7"`
	Text8  string             `bson:"text_8" json:"text_8"`
	Text9  string             `bson:"text_9" json:"text_9"`
	Text10 string             `bson:"text_10" json:"text_10"`
}