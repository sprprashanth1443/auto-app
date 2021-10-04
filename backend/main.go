package main

import (
	"github.com/auto-app/backend/user"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	router := mux.NewRouter()
	user.RegisterRoutes(router)

	log.Println("started...")
	log.Panicln(http.ListenAndServe("127.0.0.1:3000", router))
}
