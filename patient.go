package main

import "github.com/gocql/gocql"

type Patient struct {
	PatientUUID     gocql.UUID
	Age             int    `json:"age"`
	Gender          string `json:"gender"`
	InsuranceNumber string `json:"insuranceNumber"`
	Name            string `json:"name"`
}

type Patients []Patient

func PatientMapJSON(p Patient) Patient {
	
	return p
}