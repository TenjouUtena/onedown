package backend

import (
	"github.com/TenjouUtena/onedown/backend/configuration"
	"os"
	"strconv"
	"strings"

	"github.com/TenjouUtena/onedown/backend/credentials"
	"github.com/TenjouUtena/onedown/backend/session"
	"github.com/TenjouUtena/onedown/backend/users"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogger() {
	// initiate logger
	if configuration.Get().Logfile != "" {
		logfile, err := os.Open(configuration.Get().Logfile)
		if err == nil {
			log.Logger = zerolog.New(logfile).With().Timestamp().Logger()
		}
	}
	switch strings.ToUpper(configuration.Get().LogLevel) {
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

	//Load Google Oauth Credentials
	credentials.LoadCredentials(configuration.Get().CredentialFile)
	initLogger()
	log.Info().Msg("Initializing OneDown server...")

	// set up session daemon
	go session.InitDaemon(session.SessionDaemon)

	//cassandraSession := cassandra.session
	//defer cassandraSession.Close()

	// Init gin server
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

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

	users.ConfigureHandlers(router)
	credentials.ConfigureHandlers(router)

	err := router.Run(":" + strconv.Itoa(configuration.Get().Port))

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize OneDown server!")
	}
}
