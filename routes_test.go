package main

import (
	"bytes"
	"fmt"
	"github.com/gocql/gocql"
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
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check if the response's uuid is correct (Expected value).
	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct patientUUID. \n The returned message is: \n %v", rec.Body.String())
	}
}

func TestUserCreateHandler(t *testing.T) {
	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get current count of patients
	numUsersBefore := session.Query("SELECT * FROM users").Iter().NumRows()

	// Make the reader using the json string
	jsonStringReader := strings.NewReader(`{
																				  "username": "tester@test.net",
																				  "password": "test",
																				  "role": "patient",
																				  "name": "Tester"
																				}`)

	// Create the request with json as body
	req, err := http.NewRequest("POST", "/users", jsonStringReader)

	// Check if any errors occured when creating the new request
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(UserCreate)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	// Check if the number of patients had changed
	numUsersAfter := session.Query("SELECT * FROM users").Iter().NumRows()
	if numUsersAfter != (numUsersBefore + 1) {
		t.Errorf("The number of users did not change")
	}
}

func TestUserGetHandler(t *testing.T) {
	// Variables used for storing the patient
	var userUUID gocql.UUID
	var username string
	var role string
	var name string

	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get the first patient in the database
	session.Query("SELECT * FROM users").Consistency(gocql.One).Scan(&username, &name,
		&role, nil, nil, &userUUID)

	var buff bytes.Buffer
	buff.WriteString("/user/useruuid/")
	buff.WriteString(userUUID.String())
	endpoint := buff.String()
	fmt.Println(endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)

	// Check if any errors occured when creating the new request
	if err != nil {
		t.Fatal(err)
	}

	// Must manually set the endpoint URI for some unknown reason.
	req.RequestURI = endpoint

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(UserGet)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check if the response's uuid is correct (Expected value).
	if !strings.Contains(rec.Body.String(), (`"userUUID":"` + userUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct userUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"role":"` + role + `"`)) {
		t.Errorf("The response message did not contain the correct username. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"name":"` + name + `"`)) {
		t.Errorf("The response message did not contain the correct name. \n The returned message is: \n %v", rec.Body.String())
	}
	e := session.Query("DELETE FROM users where username = ?", username).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestDoctorCreateHandler(t *testing.T) {
	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get current count of patients
	numDoctorsBefore := session.Query("SELECT * FROM doctors").Iter().NumRows()

	// Make the reader using the json string
	jsonStringReader := strings.NewReader(`{
																				  "name": "Doctor Name",
																				  "phoneNumber": "111-333-2222",
																				  "primaryFacility": "address",
																				  "primarySpeciality": "Specialty",
																				  "gender": "Male"
																				}`)

	// Create the request with json as body
	req, err := http.NewRequest("POST", "/doctors", jsonStringReader)

	// Check if any errors occured when creating the new request
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(DoctorCreate)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	// Check if the number of patients had changed
	numDoctorsAfter := session.Query("SELECT * FROM doctors").Iter().NumRows()
	if numDoctorsAfter != numDoctorsBefore+1 {
		t.Errorf("The number of doctors did not change")
	}
}

func TestDoctorGetHandler(t *testing.T) {
	// Variables used for storing the patient
	var doctorUUID gocql.UUID
	var name string
	var phone string
	var primaryFacility string
	var primarySpeciality string
	var gender string

	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get the first patient in the database
	session.Query("SELECT * FROM doctors").Consistency(gocql.One).Scan(&doctorUUID, &gender,
		&name, &phone, &primaryFacility, &primarySpeciality)

	var buff bytes.Buffer
	buff.WriteString("/doctors/doctoruuid/")
	buff.WriteString(doctorUUID.String())
	endpoint := buff.String()
	fmt.Println(endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)

	// Check if any errors occured when creating the new request
	if err != nil {
		t.Fatal(err)
	}

	// Must manually set the endpoint URI for some unknown reason.
	req.RequestURI = endpoint

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(DoctorGet)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check if the response's uuid is correct (Expected value).
	if !strings.Contains(rec.Body.String(), (`"doctorUUID":"` + doctorUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct doctorUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"gender":"` + gender + `"`)) {
		t.Errorf("The response message did not contain the correct gender. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"name":"` + name + `"`)) {
		t.Errorf("The response message did not contain the correct name. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"phoneNumber":"` + phone + `"`)) {
		t.Errorf("The response message did not contain the correct name. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"primaryFacility":"` + primaryFacility + `"`)) {
		t.Errorf("The response message did not contain the correct facility. \n The returned message is: \n %v", rec.Body.String())
	}

	e := session.Query("DELETE FROM doctors where doctoruuid = ?", doctorUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestDoctorListGetHandler(t *testing.T) {
	var doctorUUID1 gocql.UUID
	var name1 string
	var phone1 string
	var primaryFacility1 string
	var primarySpeciality1 string
	var gender1 string
	var doctorUUID2 gocql.UUID
	var name2 string
	var phone2 string
	var primaryFacility2 string
	var primarySpeciality2 string
	var gender2 string

	var err error

	doctorUUID1, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name1 = "Test Doctor"
	phone1 = "123-456-7890"
	primaryFacility1 = "FakeAddress1"
	primarySpeciality1 = "Faker1"
	gender1 = "Male"
	doctorUUID2, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name2 = "Testing Doctor"
	phone2 = "0987-654-321"
	primaryFacility2 = "Fake Address2"
	primarySpeciality2 = "Faker2"
	gender2 = "Female"

	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	session.Query("INSERT INTO doctors (doctorUUID, name, phone, primaryFacility, primarySpecialty, gender) VALUES (?,?,?,?,?,?)",
		doctorUUID1, name1, phone1, primaryFacility1, primarySpeciality1, gender1).Exec()

	session.Query("INSERT INTO doctors (doctorUUID, name, phone, primaryFacility, primarySpecialty, gender) VALUES (?,?,?,?,?,?)",
		doctorUUID2, name2, phone2, primaryFacility2, primarySpeciality2, gender2).Exec()

	endpoint := "/doctors"
	fmt.Println(endpoint)

	req, err := http.NewRequest("GET", endpoint, nil)

	// Check if any errors occured when creating the new request
	if err != nil {
		t.Fatal(err)
	}

	// Must manually set the endpoint URI for some unknown reason.
	req.RequestURI = endpoint

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(DoctorListGet)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check if the response's uuid is correct (Expected value).
	if !strings.Contains(rec.Body.String(), (`"doctorUUID":"` + doctorUUID1.String() + `"`)) {
		t.Errorf("The response message did not contain the correct doctorUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"gender":"` + gender1 + `"`)) {
		t.Errorf("The response message did not contain the correct gender. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"name":"` + name1 + `"`)) {
		t.Errorf("The response message did not contain the correct name. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"phoneNumber":"` + phone1 + `"`)) {
		t.Errorf("The response message did not contain the correct name. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"primaryFacility":"` + primaryFacility1 + `"`)) {
		t.Errorf("The response message did not contain the correct facility. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"doctorUUID":"` + doctorUUID2.String() + `"`)) {
		t.Errorf("The response message did not contain the correct doctorUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"gender":"` + gender2 + `"`)) {
		t.Errorf("The response message did not contain the correct gender. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"name":"` + name2 + `"`)) {
		t.Errorf("The response message did not contain the correct name. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"phoneNumber":"` + phone2 + `"`)) {
		t.Errorf("The response message did not contain the correct name. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"primaryFacility":"` + primaryFacility2 + `"`)) {
		t.Errorf("The response message did not contain the correct facility. \n The returned message is: \n %v", rec.Body.String())
	}

	e := session.Query("DELETE FROM doctors where doctoruuid = ?", doctorUUID1).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM doctors where doctoruuid = ?", doctorUUID2).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestFutureAppointmentCreateHandler(t *testing.T) {
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

func TestFutureAppointmentGetHandler(t *testing.T) {
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
		&appointmentUUID, &dateScheduled, &doctorUUID, &notes, &patientUUID)

	fmt.Println("Querying AppointmentUUID: ", appointmentUUID.String())

	var buff bytes.Buffer
	buff.WriteString("/futureappointments/appointmentuuid/")
	buff.WriteString(appointmentUUID.String())
	endpoint := buff.String()

	fmt.Println("Using endpoint: ", endpoint)

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
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusOK)
	}

	// Check the body for patientUUID
	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct patientUUID. \nMessage: %v \nExpected:%v", rec.Body.String(), patientUUID.String())
	}

}

func TestCompletedAppointmentCreateHandler(t *testing.T) {
	var patientUUID gocql.UUID

	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	patientErr := session.Query("SELECT * FROM patients").Consistency(gocql.One).Scan(&patientUUID,
		nil, nil, nil, nil, nil, nil, nil, nil, nil)

	if patientErr != nil {
		t.Fatal(patientErr)
	}

	appointmentUUID := "1cf1dca9-4a4a-4f47-8201-401bbe0fb927"

	var bb bytes.Buffer
	bb.WriteString(`{"appointmentUUID":"`)
	bb.WriteString(appointmentUUID)
	bb.WriteString(`", "patientUUID":"`)
	bb.WriteString(patientUUID.String())
	bb.WriteString(`","doctorUUID": "1cf1dca9-4a4a-4f47-8201-401bbe0fb927",
									"dateVisited":1099,
									"breathingRate":10,
									"heartRate":80,
									"bloodOxygenLevel":56,
									"bloodPressure":129,
									"notes": "Test notes"
									}`)

	entry := bb.String()

	// Make the reader using this json string
	jsonStringReader := strings.NewReader(entry)

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
}

func TestCompletedAppointmentGetHandler(t *testing.T) {
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
	buff.WriteString("/completedappointments/appointmentuuid/")
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
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusOK)
	}

	// check the body of the returned message
	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct patientUUID. \nMessage: %v \nExpected:%v", rec.Body.String(), patientUUID.String())
	}

	e := session.Query("delete from completedappointments where appointmentuuid = ?", appointmentUUID).Exec()
	if e != nil {
		t.Fatal(e)
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
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusOK)
	}

	// There must be at least 1 entry because it was created in previous tests.
	// But we don't have any appointment UUIDs to check with, so just make sure that
	// We have at least the doctor for now.
	if !strings.Contains(rec.Body.String(), "1cf1dca9-4a4a-4f47-8201-401bbe0fb927") {
		t.Errorf("The response message did not contain the correct doctorUUID. \nMessage: %v \nExpected:%v", rec.Body.String(), "1cf1dca9-4a4a-4f47-8201-401bbe0fb927")
	}
}

func TestPatientGetByDoctorHandler(t *testing.T) {
	var patientUUID gocql.UUID
	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	patientErr := session.Query("SELECT * FROM patients").Consistency(gocql.One).Scan(&patientUUID,
		nil, nil, nil, nil, nil, nil, nil, nil, nil)

	if patientErr != nil {
		t.Fatal(patientErr)
	}

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
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusOK)
	}
	// Need to check db for uuids to check in the body
	if !strings.Contains(rec.Body.String(), patientUUID.String()) {
		t.Errorf("The response message did not contain the correct doctorUUID. \nMessage: %v \nExpected:%v", rec.Body.String(), patientUUID.String())
	}
}

func TestDeleteFutureAppointmentHandler(t *testing.T) {
	// Variables for Appointments
	var appointmentUUID gocql.UUID
	var patientUUID gocql.UUID

	// Connect to the database first.
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get a patient to create with
	patientErr := session.Query("SELECT * FROM patients").Consistency(gocql.One).Scan(&patientUUID,
		nil, nil, nil, nil, nil, nil, nil, nil, nil)

	if patientErr != nil {
		t.Fatal(patientErr)
	}

	// Create the appointment json
	var bb bytes.Buffer
	bb.WriteString(`{"patientUUID":"`)
	bb.WriteString(patientUUID.String())
	bb.WriteString(`","doctorUUID": "1cf1dca9-4a4a-4f47-8201-401bbe0fb927",
          "dateScheduled":1000, "notes": "Test notes"}`)
	entry := bb.String()

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

	// Get the current number of appointments
	numAppointments := session.Query("SELECT * FROM futureAppointments").Iter().NumRows()

	// Fetch the new appointemnt
	session.Query("SELECT * FROM futureAppointments").Consistency(gocql.One).Scan(
		&appointmentUUID, nil, nil, nil, nil)

	bb.WriteString("/futureappointments/appointmentuuid/")
	bb.WriteString(appointmentUUID.String())
	endpoint = bb.String()

	req, err = http.NewRequest("DELETE", endpoint, jsonStringReader)
	if err != nil {
		t.Fatal(err)
	}

	req.RequestURI = endpoint

	rec = httptest.NewRecorder()
	handler = http.HandlerFunc(FutureAppointmentDelete)
	handler.ServeHTTP(rec, req)

	status = rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusOK)
	}

	numAppointments2 := session.Query("SELECT * FROM futureAppointments").Iter().NumRows()

	if numAppointments2+1 != numAppointments {
		t.Errorf("The number of appointments before is %v and the current number of appointments is %v", numAppointments, numAppointments2)
	}
}
