package main

import "github.com/gocql/gocql"

type Doctor struct {
	DoctorUUID       gocql.UUID `json:"doctorUUID,omitempty"`
	Name             string     `json:"name"`
    Phone            string     `json:"phoneNumber"`
	PrimaryFacility  string     `json:"primaryFacility,omitempty"`
	PrimarySpecialty string     `json:"primarySpecialty,omitempty"`
	Gender           string     `json:"gender"`
}
