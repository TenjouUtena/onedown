package cassandra

import (
	"github.com/gocql/gocql"

	"log"
)

var Session *gocql.Session

func init() {
	var err error

	cluster := gocql.NewCluster("127.0.0.1")

	session, err := cluster.CreateSession()
	if err != nil {
		log.Println(err)
		return
	}

	err = session.Query("CREATE KEYSPACE IF NOT EXISTS onedown WITH REPLICATION = {'class' : 'SimpleStrategy','replication_factor':1};").Exec()
	if err != nil {
		log.Println(err)
		return
	}
}
