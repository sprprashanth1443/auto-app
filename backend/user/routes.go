package user

import (
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc(getUser, getUserHandler()).Methods("GET")
	router.HandleFunc(setUser, setUserHandler()).Methods("POST")

	router.HandleFunc(addTwitter,addTwitterHandler()).Methods("POST")


}
