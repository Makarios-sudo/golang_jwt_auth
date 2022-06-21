package database

import (
	"banknote-tracker-auth/models"
)
// Dummy user database
var Users = map[string]models.User{
	"micheal": {
		Username: "micheal",
		Password: "123455",
		Group:    1,
		Role:     "admin",
	},

	"james": {
		Username: "james",
		Password: "678910",
		Group:    2,
		Role:     "user",
	},
}
