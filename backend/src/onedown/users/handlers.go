package users

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//ConfigureHandlers Setup users endpoints
func ConfigureHandlers(router *gin.Engine) {
	router.POST("/users/new", postHandler)
}

//Post Users post handler
func postHandler(c *gin.Context) {
	user, err := jsonToUser(c.Request)

	if err != nil {
		c.JSON(500, gin.H{"Error": err})
		return
	}

	fmt.Println("Creating new user")

	user, err = insertNewUser(user)

	if err != nil {
		c.JSON(500, gin.H{"Error": err})
		return
	}

	c.JSON(200, gin.H{"ID": user.ID})
}
