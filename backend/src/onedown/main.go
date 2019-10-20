package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/TenjouUtena/onedown/backend/src/onedown/session"
	"github.com/gin-contrib/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func initLogger(cfg *Configuration) {
	// initiate logger
	if cfg.Logfile != "" {
		logfile, err := os.Open(cfg.Logfile)
		if err == nil {
			log.Logger = zerolog.New(logfile).With().Timestamp().Logger()
		}
	}
	switch strings.ToUpper(cfg.LogLevel) {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "ERROR":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "FATAL":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "PANIC":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

}

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

	initLogger(&cfg)
	log.Info().Msg("Iniitializing OneDown server...")

	// set up session daemon
	session.InitDaemon(session.SessionDaemon)

	r := gin.Default()

	r.Use(cors.Default()) // Needed to allow all API origins.
	r.Use(logger.SetLogger())
	// set up session routes
	session.InitSessionRoutes(r)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "PONG!")
	})
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
	err = r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize OneDown server!")
	}
}
