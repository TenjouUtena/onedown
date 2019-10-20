package users

import (
	"encoding/json"
	"net/http"
)

//JSONToUser Fills a user struct with json data from an httprequest
func JSONToUser(request *http.Request) (User, error) {
	var user User

	decoder := json.NewDecoder(request.Body)

	err := decoder.Decode(&user)

	return user, err
}
