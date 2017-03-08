package main

import "github.com/gocql/gocql"

type Document struct {
	DocumentUUID gocql.UUID `json:"documentUUID,omitempty"`
	PatientUUID  gocql.UUID `json:"patientUUID"`
	Filename     string     `json:"filename"`
	Content      string     `json:"content"`
	DateUploaded int        `json:"dateUploaded"`
}
