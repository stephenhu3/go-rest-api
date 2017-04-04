package main

import "github.com/gocql/gocql"

type Notification struct {
	Date      			int        `json:"date"`
	Messsage        	string     `json:"message"`
	ReceiverUUID    	gocql.UUID `json:"receiverUUID"`
	SenderName      	string     `json:"senderName"`
	SenderUUID      	gocql.UUID `json:"senderUUID"`
}

type Notifications []Notification