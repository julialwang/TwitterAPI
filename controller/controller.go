package controller

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"twitter-feed/config/db"
	"twitter-feed/model"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &user)
	var res model.ResponseResult

	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	collection, err := db.GetDBCollection()
	if err != nil {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
	var result model.User
	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		if err.Error() == "mongo: no documents in result" {
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
			res.Result = "Registration successful"
			json.NewEncoder(w).Encode(res)
			return
		}

		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}

	res.Result = "Username already exists"
	json.NewEncoder(w).Encode(res)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}

	collection, err := db.GetDBCollection()

	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	var res model.ResponseResult

	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		res.Error = "Invalid password"
		json.NewEncoder(w).Encode(res)
		return
	}

	if result.LoggedIn {
		res.Error = "User is already logged in"
		json.NewEncoder(w).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  result.Username,
		"firstname": result.FirstName,
		"lastname":  result.LastName,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		res.Error = "Error while generating token, try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	result.Token = tokenString
	result.Password = ""

	_, err = collection.UpdateOne(
		context.TODO(),
		bson.D{{"username", user.Username}},
		bson.D{{"$set",
			bson.D{
				{"logged-in", true},
				{"timestamp", primitive.Timestamp{T: uint32(time.Now().Unix())}},
			},
		}},
	)
	if err != nil {
		log.Fatal(err)
	}

	res.Result = "Login successful"

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		return
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	collection, err := db.GetDBCollection()

	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	var res model.ResponseResult

	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}
	if !result.LoggedIn {
		res.Error = "User is not logged in"
		json.NewEncoder(w).Encode(res)
	} else {
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.D{{"username", user.Username}},
			bson.D{{"$set",
				bson.D{
					{"logged-in", false},
				},
			}},
		)

		if err != nil {
			log.Fatal(err)
		}
		res.Result = "Logout successful"
		json.NewEncoder(w).Encode(res)
	}
	json.NewEncoder(w).Encode(result)
	return
}

func FollowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		log.Fatal(err)
	}
	collection, err := db.GetDBCollection()

	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	var res model.ResponseResult

	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&result)

	if err != nil {
		res.Error = "Invalid username"
		log.Fatal(err)
		json.NewEncoder(w).Encode(res)
		return
	}
	if !result.LoggedIn {
		res.Error = "User is not logged in"
		json.NewEncoder(w).Encode(res)
	} else {
		if result.Followings == nil {
			_, err = collection.UpdateOne(
				context.TODO(),
				bson.D{{"username", result.Username}},
				bson.D{{"$set",
					bson.D{
						{"followings", make([]string, 0)},
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
			res.Error = "Cannot follow; username invalid"
			json.NewEncoder(w).Encode(res)
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		res.Result = "Successfully followed"
		json.NewEncoder(w).Encode(res)
	}
	json.NewEncoder(w).Encode(result)
	return
}
