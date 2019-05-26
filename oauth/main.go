package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var config *oauth2.Config
var state = "romanreigns"

func init() {
	config = &oauth2.Config{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("SECRETKEY"),
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://127.0.0.1:8000/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/youtube", "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	}
}
func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/redirect", redirect)
	http.HandleFunc("/callback", callback)
	fmt.Println("Starting server...")
	http.ListenAndServe(":8000", nil)
}
func index(w http.ResponseWriter, r *http.Request) {
	var html = "<a href=\"redirect\">Sign in with google</a>"
	fmt.Fprint(w, html)
}
func redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, config.AuthCodeURL(state, oauth2.ApprovalForce, oauth2.AccessTypeOffline), http.StatusTemporaryRedirect)
}
func callback(w http.ResponseWriter, r *http.Request) {
	token := &oauth2.Token{}
	data, _ := ioutil.ReadFile("oauth/token.json")
	err := json.Unmarshal(data, token)
	if err != nil {
		fmt.Printf(err.Error())
	}
	err = requestAPI(token, w)
	if err != nil {
		token = getFreshToken(r, token)
		saveToken(token)
		requestAPI(token, w)
	}

}
func saveToken(token *oauth2.Token) {
	fmt.Println(token.RefreshToken)
	data, _ := json.MarshalIndent(token, "", "	")
	ioutil.WriteFile("oauth/token.json", data, 0777)
}
func requestAPI(token *oauth2.Token, w http.ResponseWriter) error {
	if token != nil {
		client := config.Client(oauth2.NoContext, token)
		resp, _ := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if resp == nil {
			fmt.Println("Could not get response")
			return fmt.Errorf("Error")
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Fprint(w, bodyString)
		}
		fmt.Println(resp.Status, "token saved")
	}
	return nil
}
func getFreshToken(r *http.Request, token *oauth2.Token) *oauth2.Token {
	returnstate := r.FormValue("state")
	if returnstate != state {
		fmt.Println("Invalid return state")
		return token
	}
	code := r.FormValue("code")
	token, _ = config.Exchange(oauth2.NoContext, code, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	return token
}
