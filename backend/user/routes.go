package user

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc(getUser, getUserHandler()).Methods("GET")
	router.HandleFunc(setUser, setUserHandler()).Methods("POST")

	router.HandleFunc(addTwitter, addTwitterHandler()).Methods("GET")
	router.HandleFunc(authTwitter, authTwitterHandler()).Methods("GET")
	router.HandleFunc(callBackURL, callBackURLHandler()).Methods("GET")

	router.HandleFunc(addTweetForIncentive, addTweetForIncentiveHandler()).Methods("POST")
	router.HandleFunc(submitTweetForIncentive, submitTweetForIncentiveHandler()).Methods("POST")
}
