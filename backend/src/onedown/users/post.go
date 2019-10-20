package users

import (
	"fmt"

	"github.com/TenjouUtena/onedown/backend/src/onedown/cassandra"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
)

//Post Users post handler
func Post(c *gin.Context) {
	var gocqlUUID gocql.UUID

	user, err := JSONToUser(c.Request)

	if err != nil {
		c.JSON(500, gin.H{"Error": err})
		return
	}

	fmt.Println("Creating new user")

	gocqlUUID = gocql.TimeUUID()

	err = cassandra.Session.Query("INSERT INTO users (ID, email) VALUES(?, ?);", gocqlUUID, user.Email).Exec()

	if err != nil {
		c.JSON(500, gin.H{"Error": err})
		return
	}

	c.JSON(200, gin.H{"ID": gocqlUUID})
}
