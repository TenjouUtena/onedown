package main

import (
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"os"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gp := os.Getenv("GOPATH")
	ap := path.Join(gp, "puzzles")

	r := gin.Default()

	r.Use(cors.Default()) // Needed to allow all API origins.
	r.GET("/puzzle/:puzid/get", func(c *gin.Context) {
		finalPath := path.Join(ap, c.Param("puzid")+".puz")
		puzFile, err := os.Open(finalPath)
		if err != nil {
			c.JSON(500, gin.H{"Error": err})
		} else {
			puzz, err := puzzle.ReadPuzfile(puzFile)
			if err != nil {
				c.JSON(500, gin.H{"Error": err})
			} else {
				c.JSON(200, puzz.ToPuzzle())
			}
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
