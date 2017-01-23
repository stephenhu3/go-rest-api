package main

import "github.com/gocql/gocql"

type Patient struct {
	PatientUUID     gocql.UUID `json:"patientUUID"`
	Address 		string		`json:"address"`
	BloodType		string		`json:"bloodType"`
	DateOfBirth 	int			`json:"dateOfBirth"`
	EmergencyContact string		`json:"emergencyContact"`
	Gender 			string		`json:"gender"`
	MedicalNumber	string		`json:"medicalNumber"`
	Name			string		`json:"name"`
	Notes 			string		`json:"notes"`
	Phone 			string		`json:"phoneNumber"`
}


type Patients []Patient
