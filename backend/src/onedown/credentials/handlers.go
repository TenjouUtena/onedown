package credentials

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/TenjouUtena/onedown/backend/src/onedown/users"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

//ConfigureHandlers Setup credentials api endpoints
func ConfigureHandlers(router *gin.Engine) {
	router.GET("/login", loginHandler)
	router.GET("/oauth2callback", authHandler)
}

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func authHandler(c *gin.Context) {
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	if retrievedState != c.Query("state") {
		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
		return
	}

	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	client := conf.Client(oauth2.NoContext, tok)

	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer email.Body.Close()
	data, _ := ioutil.ReadAll(email.Body)
	log.Println("Email body: ", string(data))
	u := users.GoogleUser{}
	if err = json.Unmarshal(data, &u); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	
	var user users.User
	var found bool

	if user, found = users.GetUserByEmail(u.Email); !found {
		user, err = users.InsertNewUserByGoogleUser(u)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}

	c.JSON(200, gin.H{"ID": user.ID, "Found": found})
}

func loginHandler(c *gin.Context) {
	state := randToken()
	session := sessions.Default(c)
	session.Set("state", state)
	session.Save()
	c.Writer.Write([]byte("<html><body> <a href='" + getLoginURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
}
