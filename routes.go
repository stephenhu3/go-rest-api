package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Authenticate",
		"POST",
		"/login",
		Authenticate,
	},
	Route{
		"PatientCreate",
		"POST",
		"/patients",
		PatientCreate,
	},
	Route{
		"PatientGet",
		"GET",
		"/patients/patientuuid/{patientuuid}",
		PatientGet,
	},
	Route{
		"PatientGetByDoctor",
		"GET",
		"/patients/doctoruuid/{doctoruuid}",
		PatientGetByDoctor,
	},
	Route{
		"FutureAppointmentCreate",
		"POST",
		"/futureappointments",
		FutureAppointmentCreate,
	},
	Route{
		"FutureAppointmentGet",
		"GET",
		"/futureappointments/search",
		FutureAppointmentGet,
	},
	Route{
		"CompletedAppointmentCreate",
		"POST",
		"/completedappointments",
		CompletedAppointmentCreate,
	},
	Route{
		"CompletedAppointmentGet",
		"GET",
		"/completedappointments/search",
		CompletedAppointmentGet,
	},
	Route{
		"AppointmentGetByDoctor",
		"GET",
		"/appointments/doctoruuid/{doctoruuid}",
		AppointmentGetByDoctor,
	},
	Route{
		"DoctorCreate",
		"POST",
		"/doctors",
		DoctorCreate,
	},
	Route{
		"DoctorGet",
		"GET",
		"/doctor/doctoruuid/{doctoruuid}",
		DoctorGet,
	},
}
