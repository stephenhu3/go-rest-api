package main

import "github.com/gocql/gocql"

type Prescription struct {
	PatientUUID 		gocql.UUID `json:"patientUUID"`
	PrescriptionUUID 	gocql.UUID `json:"PrescriptionUUID"`
	DoctorUUID 			gocql.UUID `json:"DoctorUUID"`
	DoctorName 			string     `json:"doctor,omitempty"`
	Drug 				string     `json:"drug"`
	StartDate 			int        `json:"startDate"`
	EndDate 			int        `json:"endDate"`
	Instructions 		string     `json:"instructions,omitempty"`
}

type Prescriptions []Prescription
