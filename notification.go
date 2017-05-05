package main

import "github.com/gocql/gocql"

type Notification struct {
	DateCreated  int        `json:"dateCreated,omitempty"`
	Messsage     string     `json:"message"`
	ReceiverUUID gocql.UUID `json:"receiverUUID"`
	SenderName   string     `json:"senderName"`
	SenderUUID   gocql.UUID `json:"senderUUID"`
}

type Notifications []Notification
