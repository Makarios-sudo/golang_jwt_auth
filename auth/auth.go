package auth

import (
	"errors"
	"github.com/go-ldap/ldap/v3"
	"log"
)

const ldapURL = "ldaps://ldap.example.com:636"

// AuthenticateConn authenticates the user login for active directory server 
// with the given username and password provided
func AuthenticateConn(username, password string) (*ldap.SearchResult, error) {
	conn, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Println("Dialurl:", err)
		return nil, errors.New("something went wrong")
	}
	defer conn.Close()

	err = conn.Bind(username, password)
	if err != nil {
		log.Println("Bind:", err)
		return nil, errors.New("something went wrong")
	}

	searchRequest := &ldap.SearchRequest{
		BaseDN:       "",
		Scope:        ldap.ScopeWholeSubtree,
		DerefAliases: ldap.NeverDerefAliases,
		SizeLimit:    0,
		TimeLimit:    0,
		TypesOnly:    false,
		Filter:       "",
		Attributes:   []string{},
		Controls:     nil,
	}

	result, err := conn.Search(searchRequest)
	if err != nil {
		log.Println("search:", err)
		return nil, errors.New("something went wrong")
	}

	return result, nil
}
