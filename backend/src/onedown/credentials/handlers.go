package credentials

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/gin-gonic/gin"
)

func getLoginURL(state string) string {
	return conf.AuthCodeURL(state)
}

func randToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}

func LoginHandler(c *gin.Context) {
	state := randToken()

	c.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + getLoginURL(state) + "'><button>Login with Google!</button> </a> </body></html>"))
}
