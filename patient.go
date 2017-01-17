package main

import "github.com/gocql/gocql"

type Patient struct {
	PatientUUID     gocql.UUID `json:"patientUUID"`
	Name			string		`json:"name"`
	Age 			int			`json:"age"`
	Gender 			string		`json:"gender"`
	InsuranceNumber	string		`json:"medicalNumber"`
	DateOfBirth 	string		`json:"dateOfBirth"`
	DateOfDeath 	string		`json:"Details"`
	Ethnicity 		string		`json:"ethnicity"`
	
	Address 		string		`json:"address"`
	PhoneNumber 	string		`json:"phoneNumber"`
	Notes 			string		`json:"notes"`
	EmerContact 	EmergencyContact `json:"emergencyContact"`
}

type EmergencyContact struct {
	Name 			string		`json:"name"`
	PhoneNumber 	string		`json:"phoneNumber"`
	Address 		string		`json:"address"`
	Relationship 	string		`json:"relationship"`
}


type Patients []Patient
