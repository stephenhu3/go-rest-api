package main

import "github.com/gocql/gocql"

type Patient struct {
	PatientUUID     gocql.UUID `json:"patientUUID"`
	Age             int        `json:"age"`
	Gender          string     `json:"gender"`
	InsuranceNumber string     `json:"insuranceNumber"`
	Name            string     `json:"name"`
}

type PatientListEntry struct {
	PatientUUID     gocql.UUID 		`json:"patientUUID"`
	Details         PatientDetails	`json:"Details"`
}

type PatientDetails struct {
	Name			string		`json:"name"`
	Age 			int			`json:"age"`
	Gender 			string		`json:"gender"`
	InsuranceNumber string		`json:"insuranceNumber"`
	DateOfBirth 	string		`json:"DOB"`
	Address 		string		`json:"address"`
	PhoneNumber 	string		`json:"phoneNumber"`
}

type PatientByList []PatientListEntry


type Patients []Patient
