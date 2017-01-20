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
