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
	// Todo sample routes
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		TodoIndex,
	},
	Route{
		"TodoCreate",
		"POST",
		"/todos",
		TodoCreate,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		TodoShow,
	},
	// End todo samples
	Route{
		"PatientCreate",
		"POST",
		"/patients",
		PatientCreate,
	},
	Route{
		"PatientGet",
		"GET",
		"/patients/search",
		PatientGet,
	},
	Route{
		"PatientGetList",
		"GET",
		"/patients/searchList",
		PatientGetList,
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
}
