package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/auto-app/backend/models"
	"github.com/auto-app/backend/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"sort"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
)

const (
	getUser     = "/user/{address}"
	setUser     = "/user"
	addTwitter  = "/addTwitter"
	callBackURL = "/twitter/callback"
	authTwitter = "/twitter"
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
		type ProviderIndex struct {
			Providers    []string
			ProvidersMap map[string]string
		}

		goth.UseProviders(
			twitter.New("", "",
				"https://www.theidentityhub.com/autonomy/authenticate/processaccountproviderresponse"),
			// If you'd like to use authenticate instead of authorize in Twitter provider, use this instead.
			// twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),
		)

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
		goth.UseProviders(
			twitter.New("", "",
				"https://www.theidentityhub.com/autonomy/authenticate/processaccountproviderresponse"),
			// If you'd like to use authenticate instead of authorize in Twitter provider, use this instead.
			// twitter.NewAuthenticate(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitter/callback"),
		)

		fmt.Println("prov", r.URL.Query().Get("provider"))

		gothic.Store = sessions.NewCookieStore([]byte(""))
	

		if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
			t, _ := template.ParseFiles("templates/success.html")
			t.Execute(w, gothUser)
		} else {
			fmt.Println("step4")
			gothic.BeginAuthHandler(w, r)
		}
	}
}

func callBackURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.URL.Query().Add("provider", "twitter")

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(w, user)

	}
}
