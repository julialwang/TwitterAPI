package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Username   string   `json:"username"`
	FirstName  string   `json:"firstname"`
	LastName   string   `json:"lastname"`
	Password   string   `json:"password"`
	Token      string   `json:"token"`
	LoggedIn   bool     `json:"logged-in" bson:"logged-in"`
	Bio        string   `json:"bio" bson:"bio"`
	Tweets     []string `json:"tweets" bson:"tweets"`
	Followings []string `json:"followings" bson:"followings"`
	Followers  []string `json:"followers" bson:"followers"`
	ToFollow   string   `json:"to-follow" bson:"to-follow"`
	NewTweet   string   `json:"new-tweet" bson:"new-tweet"`
	// separate username/following mapping
	// separate tweet table out, add dates and stuff
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

type Tweet struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Text  string             `json:"text" bson:"text"`
	Date  string             `json:"date" bson:"date"`
	Time  time.Time          `json:"time" bson:"time"`
	Owner Owner              `json:"owner" bson:"owner"`
}

func NewTweet() *Tweet {
	var t Tweet
	return &t
}
