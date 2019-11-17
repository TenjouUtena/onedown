package users

import (
	"github.com/gocql/gocql"
)

//User struct with profile data
type User struct {
	ID       gocql.UUID `json:"id"`
	Email    string     `json:"email"`
	Username string     `json:"username"`
}

//GoogleUser Google API Info
type GoogleUser struct {
	Sub           string `json:"sub"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

//GetUser Return a single User
type GetUser struct {
	User User `json:"user"`
}

//AllUsers Return an array of all users
type AllUsers struct {
	Users []User `json:"users"`
}

//NewUser New user resource ID
type NewUser struct {
	ID gocql.UUID `json:"id"`
}

//Error Return an array of error strings
type Error struct {
	Errors []string `json:"errors"`
}
