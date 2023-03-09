package auth

import "time"

type User struct {
    Username string	`xorm:"'username' unique"`
    Password string
    Salt     string
    Token    string
    TokenTime time.Time
}

type UserJSON struct {
	Username	string	`json:"username"`
	Password	string	`json:"password"`
	Token		string	`json:"token"`
}

type TokenRequest struct {
	Token    string `json:"token"`
	Username string	`json:"username"`
}