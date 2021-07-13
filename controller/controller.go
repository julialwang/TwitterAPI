package controller

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"twitter-feed/config/db"
	"twitter-feed/model"
)

var collection *mongo.Collection

// Starts the MongoDB database
func init() {
	var err error
	collection, err = db.GetDBCollection()
	if err != nil {
		log.Fatal(err)
	}
}

// RegisterHandler Registers a new user provided that the username is unique and password is valid
// Requires: username, firstname, lastname, password
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var result model.User
	var res model.ResponseResult
	var user model.User
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			if len(user.Password) < 8 || !strings.ContainsAny(user.Password, "1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 } 0") {
				res.Error = "Passwords must be longer than 8 characters and contain at least one number and letter."
				json.NewEncoder(w).Encode(res)
				return
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
			if err != nil {
				res.Error = "Error while hashing password, please try again"
				json.NewEncoder(w).Encode(res)
				return
			}
			user.Password = string(hash)
			_, err = collection.InsertOne(context.TODO(), user)
			if err != nil {
				res.Error = "Error while creating user, please try again"
				json.NewEncoder(w).Encode(res)
				return
			}
			_, err = collection.UpdateOne(
				context.TODO(),
				bson.D{{"username", user.Username}},
				bson.D{{"$set",
					bson.D{
						{"tweets", make([]string, 0)},
						{"followings", make([]string, 0)},
						{"followers", make([]string, 0)},
					},
				}},
			)
			res.Result = "Registration successful! Welcome to the team, @" + user.Username + "!"
			json.NewEncoder(w).Encode(res)
			return
		}
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	res.Result = "Username already exists, please try another :("
	json.NewEncoder(w).Encode(res)
	return
}

// LoginHandler Logs the user in with credentials if not already logged in and informs user otherwise
// Requires: username, password
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var result model.User
	var res model.ResponseResult
	var user model.User
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		res.Error = "Invalid username. Please try again!"
		json.NewEncoder(w).Encode(res)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))
	if err != nil {
		res.Error = "Invalid password. Please try again!"
		json.NewEncoder(w).Encode(res)
		return
	}
	if result.ActiveStatus {
		res.Error = "User is already logged in!"
		json.NewEncoder(w).Encode(res)
		return
	}
	result.Password = ""
	_, err = collection.UpdateOne(
		context.TODO(),
		bson.D{{"username", user.Username}},
		bson.D{{"$set",
			bson.D{
				{"active", true},
			},
		}},
	)
	if err != nil {
		log.Fatal(err)
	}
	res.Result = "Login successful. Welcome, " + result.FirstName + " " + result.LastName + "!"
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		return
	}
}

// LogoutHandler Logs the user out if not already logged out and informs user otherwise
// Requires: username, password
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var result model.User
	var res model.ResponseResult
	var user model.User
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}
	if !result.ActiveStatus {
		res.Error = "User @" + user.Username + " is already logged out."
		json.NewEncoder(w).Encode(res)
	} else {
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.D{{"username", user.Username}},
			bson.D{{"$set",
				bson.D{
					{"active", false},
				},
			}},
		)
		if err != nil {
			log.Fatal(err)
		}
		res.Result = "Logout successful! See you soon, " + result.FirstName + " " + result.LastName + "!"
		json.NewEncoder(w).Encode(res)
	}
	return
}

// FollowHandler Follows the desired user by adding their username to your "followings" and your username to their "followers"
// Requires: username, to-follow
// Handled edges: User should be logged in to follow others, and the username to follow should exist as a user in the DDB
func FollowHandler(w http.ResponseWriter, r *http.Request) {
	var result model.User
	var res model.ResponseResult
	var user model.User
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		res.Error = "Invalid username"
		log.Fatal(err)
		json.NewEncoder(w).Encode(res)
		return
	}
	if !result.ActiveStatus {
		res.Error = "You are not logged in -- Please authenticate before trying to follow users."
		json.NewEncoder(w).Encode(res)
	} else {
		if result.Followings == nil {
			_, err = collection.UpdateOne(
				context.TODO(),
				bson.D{{"username", result.Username}},
				bson.D{{"$set",
					bson.D{
						{"followings", make([]string, 0)},
						{"followers", make([]string, 0)},
					},
				}},
			)
		}
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.D{{"username", result.Username}},
			bson.D{{"$addToSet",
				bson.D{
					{"followings", user.ToFollow},
				},
			}},
		)
		err = collection.FindOne(context.TODO(), bson.D{{"username", user.ToFollow}}).Decode(&result)
		if err != nil {
			res.Error = "Cannot follow this user; The provided username is not a real user."
			json.NewEncoder(w).Encode(res)
			return
		}
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.D{{"username", user.ToFollow}},
			bson.D{{"$addToSet",
				bson.D{
					{"followers", user.Username},
				},
			}},
		)
		if err != nil {
			log.Fatal(err)
		}
		res.Result = "Successfully followed new user. Your new friend is @" + user.ToFollow + "!"
		json.NewEncoder(w).Encode(res)
	}
	return
}

// UnfollowHandler Unfollows the desired user by removing their username from your "followings" and your username from their "followers"
// Requires: username, to-follow
// Handled edges: User should be logged in to unfollow others, and the username to unfollow should be someone you're actually following
func UnfollowHandler(w http.ResponseWriter, r *http.Request) {
	var result model.User
	var res model.ResponseResult
	var user model.User
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		res.Error = "Invalid username"
		log.Fatal(err)
		json.NewEncoder(w).Encode(res)
		return
	}
	if !result.ActiveStatus {
		res.Error = "User is not logged in -- please authenticate before unfollowing"
		json.NewEncoder(w).Encode(res)
	} else {
		if result.Followings == nil {
			res.Error = "No one to unfollow -- you are not currently following anyone"
		}
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.D{{"username", user.Username}},
			bson.D{{"$pull",
				bson.D{
					{"followings", user.ToFollow},
				},
			}},
		)
		err = collection.FindOne(context.TODO(), bson.D{{"username", user.ToFollow}}).Decode(&result)
		if err != nil {
			res.Error = "Failed to unfollow @" + user.ToFollow + ", as you are were never actually following them in the first place."
			json.NewEncoder(w).Encode(res)
			return
		}
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.D{{"username", user.ToFollow}},
			bson.D{{"$pull",
				bson.D{
					{"followers", user.Username},
				},
			}},
		)
		if err != nil {
			log.Fatal(err)
		}
		res.Result = "Successfully unfollowed user @" + user.ToFollow + ". Bye!"
		json.NewEncoder(w).Encode(res)
	}
	return
}

// TweetHandler Tweets the input text to your profile, where it is saved in chronological order
// Requires: username, new-tweet
// Handled edges: User should be logged in to tweet, and the tweet should not only contain whitespace
func TweetHandler(w http.ResponseWriter, r *http.Request) {
	var result model.User
	var res model.ResponseResult
	var user model.User
	w.Header().Set("Content-Type", "application/json")
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)
	if err != nil {
		res.Error = "Invalid username"
		log.Fatal(err)
		json.NewEncoder(w).Encode(res)
		return
	}
	if !result.ActiveStatus {
		res.Error = "You are not logged in -- Please authenticate before tweeting!"
		json.NewEncoder(w).Encode(res)
	} else {
		if result.Tweets == nil {
			_, err = collection.UpdateOne(
				context.TODO(),
				bson.D{{"username", result.Username}},
				bson.D{{"$set",
					bson.D{
						{"tweets", make([]string, 0)},
					},
				}},
			)
		}
		if strings.TrimSpace(result.NewTweet) == "" {
			res.Error = "Aren't you going to say anything in your Tweet? Write something!"
			json.NewEncoder(w).Encode(res)
			return
		}
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.D{{"username", result.Username}},
			bson.D{{"$addToSet",
				bson.D{
					{"tweets", user.NewTweet},
				},
			}},
		)
		if err != nil {
			log.Fatal(err)
		}
		res.Result = "Successfully tweeted at " + string(time.Now().Format("01-02-2006 15:04:05"))
		json.NewEncoder(w).Encode(res)
	}
	return
}

// ProfileHandler Displays the profile of any user in the DDB provided that they exist
// Requires: {username} in request
func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{"username", params["username"]}}).Decode(&result)
	json.NewEncoder(w).Encode(result)
	return
}
