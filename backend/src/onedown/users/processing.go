package users

import (
	"encoding/json"
	"net/http"

	"github.com/TenjouUtena/onedown/backend/src/onedown/cassandra"
	"github.com/gocql/gocql"
)

//JSONToUser Fills a user struct with json data from an httprequest
func JSONToUser(request *http.Request) (User, error) {
	var user User

	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&user)

	return user, err
}

//GetUserByEmail Return single user by UUID. Bool is whether the user was found
func GetUserByEmail(email string) (User, bool) {
	var user User
	var found bool = false

	m := map[string]interface{}{}
	query := "SELECT id, email, username FROM users WHERE email = ? LIMIT 1"
	iterable := cassandra.Session.Query(query, email).Consistency(gocql.One).Iter()
	for iterable.MapScan(m) {
		found = true
		user = User{
			ID:       m["id"].(gocql.UUID),
			Email:    m["email"].(string),
			Username: m["username"].(string),
		}
	}

	return user, found
}
