package cassandra

import (
	"fmt"
	"github.com/TenjouUtena/onedown/backend/configuration"
	"github.com/gocql/gocql"
	"github.com/rs/zerolog/log"
)

//session Current cassandra session
var session *gocql.Session

func GetSession() *gocql.Session {
	if session == nil {
		cluster := gocql.NewCluster(configuration.Get().CassandraHost)
		cluster.Keyspace = "onedown"
		var err error
		session, err = cluster.CreateSession()
		if err != nil {
			panic(err)
		}
	}
	return session
}

func init() {
	var err error

	err = GetSession().Query("CREATE KEYSPACE IF NOT EXISTS onedown WITH REPLICATION = {'class' : 'SimpleStrategy','replication_factor':1};").Exec()
	if err != nil {
		log.Error().Err(err).Msg("Error initializing cassandra.")
		return
	}

	err = GetSession().Query("CREATE TABLE IF NOT EXISTS users (ID uuid, Email text, Username text, PRIMARY KEY(ID, email));").Exec()
	if err != nil {
		log.Error().Err(err).Msg("Error initializing cassandra.")
		return
	}

	err = GetSession().Query("CREATE TABLE IF NOT EXISTS puzzle_sessions (ID uuid, SessionData blob, PRIMARY KEY(ID));").Exec()
	if err != nil {
		log.Error().Err(err).Msg("Error initializing cassandra.")
		return
	}

	err = GetSession().Query("CREATE TABLE IF NOT EXISTS puzzle (ID uuid, PuzzleData blob, PRIMARY KEY(ID));").Exec()
	if err != nil {
		log.Error().Err(err).Msg("Error initializing cassandra.")
		return
	}

	fmt.Println("Cassandra init done")
}
