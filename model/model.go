package model

type User struct {
	Username     string   `json:"username"`
	FirstName    string   `json:"firstname"`
	LastName     string   `json:"lastname"`
	Password     string   `json:"password"`
	ActiveStatus bool     `json:"active" bson:"active"`
	Bio          string   `json:"bio" bson:"bio"`
	Tweets       []string `json:"tweets" bson:"tweets"`
	Followings   []string `json:"followings" bson:"followings"`
	Followers    []string `json:"followers" bson:"followers"`
	ToFollow     string   `json:"to-follow" bson:"to-follow"`
	NewTweet     string   `json:"new-tweet" bson:"new-tweet"`
}

type ResponseResult struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}
