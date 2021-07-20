package model

import "github.com/google/uuid"

type User struct {
	Username     string      `json:"username"`
	FirstName    string      `json:"firstname"`
	LastName     string      `json:"lastname"`
	Password     string      `json:"password"`
	ActiveStatus bool        `json:"active" bson:"active"`
	Bio          string      `json:"bio" bson:"bio"`
	Followings   []string    `json:"followings" bson:"followings"`
	Followers    []string    `json:"followers" bson:"followers"`
	Input        string      `json:"input" bson:"input"`
	TweetIDs     []uuid.UUID `json:"tweetids" bson:"tweetids"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}
