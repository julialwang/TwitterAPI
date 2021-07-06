package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	Bio            string `json:"bio" bson:"bio"`
	Tweets        *[]primitive.ObjectID `json:"tweets" bson:"tweets"`
	Followings    *[]Owner              `json:"followings" bson:"followings"`
	Followers     *[]Owner              `json:"followers" bson:"followers"`
	// Pull the array first then update ** another table
	// bson?
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type Tweet struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Text     string             `json:"text" bson:"text"`
	Date     string             `json:"date" bson:"date"`
	Time     time.Time          `json:"time" bson:"time"`
	Owner    Owner              `json:"owner" bson:"owner"`
}

func NewTweet() *Tweet {
	var t Tweet
	return &t
}