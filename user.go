package main

import "github.com/gocql/gocql"

type User struct {
	UserName         string     	`json:"userName,omitempty"`
	PassWord         string     	`json:"passWord,,omitempty"`
	UserUUID         gocql.UUID     `json:"userUUID,,omitempty"`
	Role         	 string     	`json:"role,,omitempty"`
	Name         	 string     	`json:"name,,omitempty"`
}