package users

import (
	"encoding/json"
	"net/http"

	"github.com/TenjouUtena/onedown/backend/src/onedown/cassandra"
	"github.com/gocql/gocql"
)

//InsertNewUserByGoogleUser Add new user by googleUser struct
func InsertNewUserByGoogleUser(googleUser GoogleUser) (User, error) {
	user := User{
		Email:    googleUser.Email,
		Username: googleUser.Email,
	}

	user, err := insertNewUser(user)

	return user, err
}

func jsonToUser(request *http.Request) (User, error) {
	var user User

	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&user)

	return user, err
}

func insertNewUser(user User) (User, error) {
	user.ID = gocql.TimeUUID()
	err := cassandra.Session.Query("INSERT INTO users (ID, email, username) VALUES(?, ?, ?);", user.ID, user.Email, user.Username).Exec()

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
