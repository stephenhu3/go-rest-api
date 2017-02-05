package main

import (
	"bytes"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIndexHandler(t *testing.T) {
	// Create the request
	req, err := http.NewRequest("GET", "/index", nil)
	// Check if any errors occured when creating the new request
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(Index)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check the body of the Index page. At this moment, it is
	// Supposed to be "Welcome!\n"
	expected := "Welcome!\n"
	if rec.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v, want %v", rec.Body.String(), expected)
	}
}

// This is used as a base for checking database endpoints.
func TestPatientCreateHandler(t *testing.T) {
	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get current count of patients
	numPatientsBefore := session.Query("SELECT * FROM patients").Iter().NumRows()

	// Make the reader using the json string
	jsonStringReader := strings.NewReader(`{"age": "69",
                                          "gender": "F",
                                          "name": "Kelly",
                                          "insuranceNumber": "1234567890"
                                          }`)

	// Create the request with json as body
	req, err := http.NewRequest("POST", "/patients", jsonStringReader)

	// Check if any errors occured when creating the new request
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(PatientCreate)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	// Check if the response is code 201 (Expected).
	if !strings.Contains(rec.Body.String(), `"code":201`) {
		t.Errorf("The response message did not contain Code 201. \n The returned message is: \n %v", rec.Body.String())
	}

	// Check if the number of patients had changed
	numPatientsAfter := session.Query("SELECT * FROM patients").Iter().NumRows()
	if numPatientsAfter != (numPatientsBefore + 1) {
		t.Errorf("The number of patients did not change")
	}
}

// This is used as a base for checking database endpoints.
func TestPatientGetHandler(t *testing.T) {
	// Variables used for storing the patient
	var patientUUID gocql.UUID
	var address string
	var bloodType string
	var dateOfBirth int
	var emergencyContact string
	var gender string
	var medicalNumber string
	var name string
	var notes string
	var phone string

	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get the first patient in the database
	session.Query("SELECT * FROM patients").Consistency(gocql.One).Scan(&patientUUID, &address,
		&bloodType, &dateOfBirth, &emergencyContact, &gender, &medicalNumber,
		&name, &notes, &phone)

	var buff bytes.Buffer
	buff.WriteString("/patients/patientuuid/")
	buff.WriteString(patientUUID.String())
	endpoint := buff.String()

	req, err := http.NewRequest("GET", endpoint, nil)

	// Check if any errors occured when creating the new request
	if err != nil {
		t.Fatal(err)
	}

	// Must manually set the endpoint URI for some unknown reason.
	req.RequestURI = endpoint

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(PatientGet)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	// Check if the response's uuid is correct (Expected value).
	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct patientUUID. \n The returned message is: \n %v", rec.Body.String())
	}
}

func TestFutureAppointmentCreate(t *testing.T) {
	var patientUUID gocql.UUID
	var entry string
	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get current count of appointments
	numAppointments := session.Query("SELECT * FROM futureAppointments").Iter().NumRows()

	patientErr := session.Query("SELECT * FROM patients").Consistency(gocql.One).Scan(&patientUUID,
		nil, nil, nil, nil, nil, nil, nil, nil, nil)

	if patientErr != nil {
		t.Fatal(patientErr)
	}
	var bb bytes.Buffer
	log.Printf("Patient was found")
	bb.WriteString(`{"patientUUID":"`)
	bb.WriteString(patientUUID.String())
	bb.WriteString(`","doctorUUID": "1cf1dca9-4a4a-4f47-8201-401bbe0fb927",
          "dateScheduled":1000, "notes": "Test notes"}`)
	entry = bb.String()

	// Make the reader using this json string
	jsonStringReader := strings.NewReader(entry)

	endpoint := "/futureappointments"
	req, err := http.NewRequest("POST", endpoint, jsonStringReader)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(FutureAppointmentCreate)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusNotFound)
	}

	// Check if the number of appointments actually went up
	numAppointments2 := session.Query("SELECT * FROM futureAppointments").Iter().NumRows()

	if numAppointments2 != numAppointments+1 {
		t.Errorf("Number of appointments in the database: got %v but supposed to be %v", numAppointments2, numAppointments+1)
	}
}

func TestFutureAppointmentGet(t *testing.T) {
	// Variables for Appointments
	var appointmentUUID gocql.UUID
	var patientUUID gocql.UUID
	var doctorUUID gocql.UUID
	var dateScheduled int
	var notes string

	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get the values from the first appointment found
	session.Query("SELECT * FROM futureAppointments").Consistency(gocql.One).Scan(
		&appointmentUUID, &patientUUID, &doctorUUID, &dateScheduled, &notes)

	var buff bytes.Buffer
	buff.WriteString("/futureappointments/search/?appointmentuuid=")
	buff.WriteString(appointmentUUID.String())
	endpoint := buff.String()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Manually set the endpoint in the request URI since the
	// function isn't setting it on its own for GET requests
	req.RequestURI = endpoint

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(FutureAppointmentGet)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusFound)
	}
}

func TestCompletedAppointmentCreate(t *testing.T) {
	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get current count of appointments
	numAppointments := session.Query("SELECT * FROM completedAppointments").Iter().NumRows()

	// Make the reader using this json string
	jsonStringReader := strings.NewReader((`{"patientUUID": "1cf1dca9-4a4a-4f47-8201-401bbe0fb925",
                                          "doctorUUID": "1cf1dca9-4a4a-4f47-8201-401bbe0fb927",
                                          "dateVisited":1099,
                                          "breathingRate":10,
                                          "heartRate":80,
                                          "bloodOxygenLevel":56,
                                          "bloodPressure":129,
                                          "notes": "Test notes"
                                          }`))

	endpoint := "/completedappointments"
	req, err := http.NewRequest("POST", endpoint, jsonStringReader)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(CompletedAppointmentCreate)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusNotFound)
	}

	// Check if the number of appointments actually went up
	numAppointments2 := session.Query("SELECT * FROM completedAppointments").Iter().NumRows()

	if numAppointments2 != numAppointments+1 {
		t.Errorf("Number of appointments in the database: got %v but supposed to be %v", numAppointments2, numAppointments+1)
	}
}

func TestCompletedAppointmentGet(t *testing.T) {
	// Variables for Appointments
	var appointmentUUID gocql.UUID
	var patientUUID gocql.UUID
	var doctorUUID gocql.UUID
	var dateVisited int
	var breathingRate int
	var heartRate int
	var bloodOxygenLevel int
	var bloodPressure int
	var notes string

	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get the values from the first appointment found
	err := session.Query("SELECT * FROM completedAppointments").Consistency(gocql.One).Scan(
		&appointmentUUID, &bloodOxygenLevel, &bloodPressure,
		&breathingRate, &dateVisited, &doctorUUID, &heartRate, &notes, &patientUUID)

	if err != nil {
		t.Errorf("There are no completedAppointments")
	}

	var buff bytes.Buffer
	buff.WriteString("/completedappointments/search/?appointmentuuid=")
	buff.WriteString(appointmentUUID.String())
	endpoint := buff.String()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Manually set the endpoint in the request URI since the
	// function isn't setting it on its own for GET requests
	req.RequestURI = endpoint

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(CompletedAppointmentGet)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusFound)
	}
}

func TestPatientGetByDoctorHandler1(t *testing.T) {
	endpoint := "/patients/doctoruuid/test-uuid-should-fail"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Manually set the endpoint in the request URI since the
	// function isn't setting it on its own.
	req.RequestURI = endpoint

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(PatientGetByDoctor)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusNotFound {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusNotFound)
	}
}

func TestPatientGetByDoctorHandler(t *testing.T) {
	endpoint := "/patients/doctoruuid/1cf1dca9-4a4a-4f47-8201-401bbe0fb927"
	// The doctorUUID is the same as the UUID used for doctors in the test above.
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Manually set the endpoint in the request URI since the
	// function isn't setting it on its own.
	req.RequestURI = endpoint

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(PatientGetByDoctor)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusFound)
	}
}

func TestAppointmentGetByDoctorHandler(t *testing.T) {
	endpoint := "/appointments/doctoruuid/1cf1dca9-4a4a-4f47-8201-401bbe0fb927"
	// The doctorUUID is the same as the UUID used for doctors in the test above.
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Manually set the endpoint in the request URI since the
	// function isn't setting it on its own.
	req.RequestURI = endpoint

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(AppointmentGetByDoctor)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusFound {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusFound)
	}
}
