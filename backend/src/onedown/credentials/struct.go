package credentials

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cred Credentials
var conf *oauth2.Config

//Credentials Struct that stores googleids
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

func LoadCredentials(fileName string) {
	file, err := ioutil.ReadFile(fileName) //Currently just a json file, should be stored more securely

	if err != nil {
		fmt.Printf("File error: v\n", err)
	}

	json.Unmarshal(file, &cred)

	conf = &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  "http://127.0.0.1:8080/oauth2callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}
