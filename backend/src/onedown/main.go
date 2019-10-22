package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/TenjouUtena/onedown/backend/src/onedown/cassandra"
	"github.com/TenjouUtena/onedown/backend/src/onedown/credentials"
	"github.com/TenjouUtena/onedown/backend/src/onedown/puzzle"
	"github.com/TenjouUtena/onedown/backend/src/onedown/session"
	"github.com/TenjouUtena/onedown/backend/src/onedown/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Initializing OneDown server...")

	//Load Google Oauth Credentials
	credentials.LoadCredentials(cfg.CredentialFile)

	// set up session daemon
	go session.InitDaemon(session.SessionDaemon)

	cassandraSession := cassandra.Session
	defer cassandraSession.Close()

	// Init gin server
	router := gin.Default()

	router.Use(cors.Default()) // Needed to allow all API origins.
	router.Use(logger.SetLogger(logger.Config{
		Logger: &log.Logger,
	}))

	// set up session routes
	session.InitSessionRoutes(router)

	// set up general routes
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "PONG!")
	})

	router.GET("/puzzle/:puzid/get", func(c *gin.Context) {
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
	router.POST("/users/new", users.Post)

	router.Run(":" + strconv.Itoa(cfg.Port))

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize OneDown server!")
	}
}
