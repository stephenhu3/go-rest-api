package main

import "github.com/gocql/gocql"

type User struct {
	Username 		string     `json:"username,omitempty"`
	Password 		string     `json:"password,omitempty"`
	UserUUID 		gocql.UUID `json:"userUUID,omitempty"`
	Role     		string     `json:"role,omitempty"`
	Name     		string     `json:"name,omitempty"`
	VerificationKey	string     `json:"verificationKey,omitempty"`
}

