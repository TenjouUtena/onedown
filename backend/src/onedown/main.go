package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/TenjouUtena/onedown/backend/src/onedown/session"
	"os"
	"path"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	var cfg Configuration
	cfg.GOPATH = os.Getenv("GOPATH")

	configPath := flag.String("config", path.Join(cfg.GOPATH, "configuration"), "Path to configuration files")
	flag.Parse()
	file, err := os.Open(path.Join(*configPath, "configuration.json"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// set up session daemon
	session.InitDaemon(session.SessionDaemon)

	r := gin.Default()

	r.Use(cors.Default()) // Needed to allow all API origins.
	// set up solver routes
	session.InitSolverRoutes(r)

	r.GET("/puzzle/:puzid/get", func(c *gin.Context) {
		finalPath := path.Join(cfg.PuzzleDirectory, c.Param("puzid")+".puz")
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
