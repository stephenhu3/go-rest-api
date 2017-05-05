package main

import "github.com/gocql/gocql"

type GenericAppointment struct {
	AppointmentUUID gocql.UUID `json:"appointmentUUID"`
	PatientUUID     gocql.UUID `json:"patientUUID"`
	DoctorUUID      gocql.UUID `json:"doctorUUID"`
	PatientName     string     `json:"patientName"`
	DateScheduled   int        `json:"dateScheduled"`
	DateVisited     int        `json:"dateVisited"`
	Notes           string     `json:"notes"`
}

type GenericAppointments []GenericAppointment
