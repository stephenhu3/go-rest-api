package main

import "github.com/gocql/gocql"

type Patient struct {
	PatientUUID     gocql.UUID `json:"patientUUID"`
	Age             int    `json:"age"`
	Gender          string `json:"gender"`
	InsuranceNumber string `json:"insuranceNumber"`
	Name            string `json:"name"`
}

type Patients []Patient
