package main

import "github.com/gocql/gocql"

type Patient struct {
	PatientUUID      gocql.UUID `json:"patientUUID"`
	Address          string     `json:"address,omitempty"`
	BloodType        string     `json:"bloodType,omitempty"`
	DateOfBirth      int        `json:"dateOfBirth"`
	EmergencyContact string     `json:"emergencyContact,omitempty"`
	Gender           string     `json:"gender"`
	MedicalNumber    string     `json:"medicalNumber,omitempty"`
	Name             string     `json:"name"`
	Notes            string     `json:"notes,omitempty"`
	Phone            string     `json:"phoneNumber"`
}

type Patients []Patient
