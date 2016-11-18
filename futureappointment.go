package main

import "github.com/gocql/gocql"

type FutureAppointment struct {
	AppointmentUuid gocql.UUID `json:"appointmentUuid"`
	PatientUuid     gocql.UUID `json:"patientUuid"`
	DoctorUuid      gocql.UUID `json:"doctorUuid"`
	DateScheduled   int        `json:"dateScheduled"`
	Notes           string     `json:"notes"`
}
