package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auto-app/backend/config"
	"github.com/auto-app/backend/models"
	"github.com/auto-app/backend/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
)

const (
	getUser     = "/user/{address}"
	setUser     = "/user"
	addTwitter  = "/addTwitter"
	authTwitter = "/twitter"
	callBackURL = "/twitter/callback"

	addTweetForIncentive    = "/twitter/add"
	submitTweetForIncentive = "/twitter/submit"
)

type user struct {
	Address string `json:"address"`
	Handle  string `json:"handle"`
}

type addTweet struct {
	URL     string `json:"url"`
	Address string `json:"address"`
}

type submitTweet struct {
	URL     string `json:"url"`
	Address string `json:"address"`
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
		type ProviderIndex struct {
			Providers    []string
			ProvidersMap map[string]string
		}

		goth.UseProviders(twitter.New(config.APIKey, config.APIKeySecret, config.CallBackURL))

		m := make(map[string]string)
		m["twitter"] = "Twitter"

		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		providerIndex := &ProviderIndex{
			Providers:    keys,
			ProvidersMap: m,
		}

		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(w, providerIndex)
	}
}

func authTwitterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessions.NewCookieStore([]byte(config.APIKeySecret))
		gothic.StoreInSession("twitter", "twitter", r, w)

		if _, err := gothic.CompleteUserAuth(w, r); err == nil {

			utils.WriteErrorToResponse(w, 500, "")
			return
		} else {

			gothic.BeginAuthHandler(w, r)
			return
		}
	}
}

func callBackURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us, err := gothic.CompleteUserAuth(w, r)
		if err != nil {

			utils.WriteErrorToResponse(w, 500, "error while authenticating user")
			return
		}

		address := r.URL.Query().Get("address")
		pr := r.URL.Query().Get("provider")
		q := bson.M{"address": address}
		user := models.User{}
		user.FindOne(q)
		fmt.Println(address, pr)
		if user.Address == "" {
			user = user.NewUser(address, us.NickName)
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
func addTweetForIncentiveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t addTweet
		var unmarshalErr *json.UnmarshalTypeError

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&t)
		if err != nil {
			if errors.As(err, &unmarshalErr) {
				utils.WriteErrorToResponse(w, 500, err)
				return
			} else {
				utils.WriteErrorToResponse(w, 500, err)
				return
			}
		}

		var tw models.Tweet
		q := bson.M{"url": t.URL, "added_by": t.Address}
		err = tw.FindOne(q)
		if err != nil && err.Error() != "not found" {
			utils.WriteErrorToResponse(w, 500, err.Error())
			return
		}
		if tw.AddedBy != "" {
			utils.WriteErrorToResponse(w, 500, "addTweet already existed")
			return

		} else {
			tw = tw.NewTweet(t.URL, t.Address)
			err := tw.Save()
			if err != nil {
				utils.WriteErrorToResponse(w, 500, err)
				return
			}
			utils.WriteResultToResponse(w, 200, " addTweet added successfully")
			return
		}
	}
}

func submitTweetForIncentiveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t submitTweet
		var unmarshalErr *json.UnmarshalTypeError
		baseURL := "https://api.twitter.com/2/tweets/"

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := decoder.Decode(&t)
		if err != nil {
			if errors.As(err, &unmarshalErr) {
				fmt.Println("1")
				utils.WriteErrorToResponse(w, 500, err)
				return
			} else {
				fmt.Println("2")
				utils.WriteErrorToResponse(w, 500, err)
				return
			}
		}

		var u models.User
		q := bson.M{"address": t.Address}
		err = u.FindOne(q)
		if err != nil {
			utils.WriteErrorToResponse(w, 500, err.Error())
			return
		}
		if u.Handle == "" {
			fmt.Println("4")
			utils.WriteErrorToResponse(w, 500, "Twitter account not verified")
			return
		}

		strs := strings.Split(t.URL, "/")
		baseURL = baseURL + strs[len(strs)-1] + "/liking_users?"

		http.Request{Header: {"bearer_token":}}
		res, err := http.Get(baseURL)
		if err != nil {
			fmt.Println("5")
			utils.WriteErrorToResponse(w, 500, err)
			return
		}

		body, err := ioutil.ReadAll(res.Body)

		w.Write(bodyz)
		//utils.WriteResultToResponse(w, 200, res)
		return
	}
}
