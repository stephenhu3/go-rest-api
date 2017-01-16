package main

import "github.com/gocql/gocql"

type CompletedAppointment struct {
	AppointmentUUID  gocql.UUID `json:"appointmentUUID"`
	PatientUUID      gocql.UUID `json:"patientUUID"`
	DoctorUUID       gocql.UUID `json:"doctorUUID"`
	DateVisited      int        `json:"dateVisited"`
	BreathingRate    int        `json:"breathingRate"`
	HeartRate        int        `json:"heartRate"`
	BloodOxygenLevel int        `json:"bloodOxygenLevel"`
	BloodPressure    int        `json:"bloodPressure"`
	Notes            string     `json:"notes"`
}

type CompletedAppointments []CompletedAppointment
