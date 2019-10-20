package cassandra

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

//Session Current cassandra Session
var Session *gocql.Session

func init() {
	var err error

	cluster := gocql.NewCluster("127.0.0.1")

	Session, err = cluster.CreateSession()
	if err != nil {
		log.Println(err)
		return
	}

	err = Session.Query("CREATE KEYSPACE IF NOT EXISTS onedown WITH REPLICATION = {'class' : 'SimpleStrategy','replication_factor':1};").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	cluster.Keyspace = "onedown"

	err = Session.Query("CREATE TABLE IF NOT EXISTS onedown.users (ID uuid, Email text, PRIMARY KEY(ID, email));").Exec()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Cassandra init done")
}