package cassandra

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

//Session Current cassandra session
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

	cluster.Keyspace = "onedown"

	err = session.Query("CREATE TABLE IF NOT EXISTS onedown.users (ID uuid, Email text, PRIMARY KEY(ID, email));").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Cassandra init done")
}
