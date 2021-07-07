package model

type Owner struct {
	Username    string `json:"username"`
	IsFollowing bool   `json:"is_following"`
}

func NewOwner(username string, pic string, name string, bio string) *Owner {
	return &Owner{
		Username: username,
	}
}
