package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"twitter-feed/controller"
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
	r.HandleFunc("/profile/{username}", controller.ProfileHandler).
		Methods("GET")
	r.HandleFunc("/timeline", controller.TimelineHandler).
		Methods("GET")
	r.HandleFunc("/delete", controller.DeleteHandler).
		Methods("POST")
	r.HandleFunc("/untweet", controller.UntweetHandler).
		Methods("POST")
	r.HandleFunc("/update", controller.UpdateHandler).
		Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
