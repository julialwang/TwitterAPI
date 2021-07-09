package main

import (
	"github.com/gorilla/mux"
	"twitter-feed/controller"
	//"twitter-feed/model"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register", controller.RegisterHandler).
		Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).
		Methods("POST")
	r.HandleFunc("/logout", controller.LogoutHandler).
		Methods("POST")
	r.HandleFunc("/follow", controller.FollowHandler).
		Methods("POST")
	r.HandleFunc("/unfollow", controller.UnfollowHandler).
		Methods("POST")
	r.HandleFunc("/tweet", controller.TweetHandler).
		Methods("POST")
	r.HandleFunc("/fetch", controller.FetchHandler).
		Methods("POST")
	r.HandleFunc("/profile/{username}", controller.ProfileHandler).
		Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
