package models

import "github.com/dgrijalva/jwt-go"

// LoginData hold the username and password provided by the client
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User is a dummy user for testing
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Group    int    `json:"group"`
	Role     string `json:"role"`
}

// UserNoPassword is a user with no password field
type UserNoPassword struct {
	Username string `json:"username"`
	Group    int    `json:"group"`
	Role     string `json:"role"`
}

// Claims models type claims from jwt library
type Claims struct {
	UserNoPassword
	jwt.StandardClaims
}

type Token struct {
	Token string `json:"token"`
}
