package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)


var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// TODO RANDOMIZE
	randomState = "random"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", welcomePage)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/callback", handleCallback)
	http.ListenAndServe("localhost:8080", r)

}

func welcomePage(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/login">Google Login</a></body></html>`
	fmt.Fprint(w, html)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	url := googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		fmt.Println("State is not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Printf("Could not get token %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("Cannot do get request %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Cannot parse response %s\n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Response %s = ", content.)
}
