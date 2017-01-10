package main

import "github.com/gocql/gocql"

type FutureAppointment struct {
	AppointmentUUID gocql.UUID `json:"appointmentUUID"`
	PatientUUID     gocql.UUID `json:"patientUUID"`
	DoctorUUID      gocql.UUID `json:"doctorUUID"`
	DateScheduled   int        `json:"dateScheduled"`
	Notes           string     `json:"notes"`
}

type FutureAppointments []FutureAppointment
