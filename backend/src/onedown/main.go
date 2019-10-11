
package main

import "github.com/gin-gonic/gin"
import "github.com/gin-contrib/cors"

func main() {
	r := gin.Default()

	r.Use(cors.Default())  // Needed to allow all API origins.
	r.GET("/puzzle/:puzid/get", func(c *gin.Context) {
		c.JSON(200, []map[string]int{
         map[string]int{"x":1, "y":1},
         map[string]int{"x":1, "y":2},
         map[string]int{"x":2, "y":1},
         map[string]int{"x":2, "y":2}})
  })
	r.Run() // listen and serve on 0.0.0.0:8080
}




