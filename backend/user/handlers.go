package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auto-app/backend/models"
	"github.com/auto-app/backend/utils"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

const (
	getUser    = "/user/{address}"
	setUser    = "/user"
	addTwitter = "/user/addTwitter"
)

type user struct {
	Address string `json:"address"`
	Handle  string `json:"handle"`
}

func getUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		q := bson.M{"address": vars["address"]}
		user := models.User{}
		user.FindOne(q)
		if user.Address == "" {
			utils.WriteErrorToResponse(w, 500, "No user found")
			return
		}

		utils.WriteResultToResponse(w, 200, user)
		return
	}
}

func setUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u user
		var unmarshalErr *json.UnmarshalTypeError

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&u)
		if err != nil {
			if errors.As(err, &unmarshalErr) {
				utils.WriteErrorToResponse(w, 500, err)
				return
			} else {
				utils.WriteErrorToResponse(w, 500, err)
				return

			}
		}

		q := bson.M{"address": u.Address}
		user := models.User{}
		fmt.Println(user)
		user.FindOne(q)
		if user.Address == "" {
			user = user.NewUser(u.Address, u.Handle)
			err = user.Save()
			if err != nil {
				utils.WriteErrorToResponse(w, 500, err)
				return
			}
		} else {
			utils.WriteErrorToResponse(w, 500, "User already exists")
			return
		}

		utils.WriteResultToResponse(w, 200, user)
		return

	}
}

func addTwitterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		conf := &oauth2.Config{
			ClientID:     "1214434556872232962-cG4ZeDmY9koeOQpc6q2nAbJIcJLt5j",
			ClientSecret: "nPQ6TdR9ZeoWpXWSGivoYGuByL3FP5bMgj5LtSxCz6CzU",
			Scopes:       []string{},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://api.twitter.com/oauth2/auth",
				TokenURL: "https://api.twitter.com/oauth2/token",
			},
		}


		url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
		fmt.Printf("Visit the URL for the auth dialog: %v", url)

		var code string
		if _, err := fmt.Scan(&code); err != nil {
			log.Fatal(err)
		}
		tok, err := conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}

		client := conf.Client(ctx, tok)
		client.Get("...")
	}
}
