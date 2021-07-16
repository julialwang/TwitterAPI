package model

import (
	"github.com/google/uuid"
)

type Tweet struct {
	ID   uuid.UUID `json:"id,omitempty" bson:"_id"`
	Text string    `json:"text" bson:"text"`
	Date string    `json:"date" bson:"date"`
	Time string    `json:"time" bson:"time"`
}

type TweetResp struct {
	User string `json:"user"`
	Date string `json:"date" bson:"date"`
	Time string `json:"time" bson:"time"`
	Text string `json:"text" bson:"text"`
}

type Timeline struct {
	Tweets []TweetResp `json:"tweets"`
}
