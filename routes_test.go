package main

import (
	"bytes"
	"fmt"
	"github.com/gocql/gocql"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

var testDB string = "emr"

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
	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = testDB
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

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = testDB
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
	if !strings.Contains(rec.Body.String(), (`"gender":"` + gender + `"`)) {
		t.Errorf("The response message did not contain the correct gender. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"name":"` + name + `"`)) {
		t.Errorf("The response message did not contain the correct name. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"phoneNumber":"` + phone + `"`)) {
		t.Errorf("The response message did not contain the correct phone number. \n The returned message is: \n %v", rec.Body.String())
	}
}

func TestUserCreateHandler(t *testing.T) {
	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = testDB
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

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = testDB
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
	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get current count of patients
	numDoctorsBefore := session.Query("SELECT * FROM doctors").Iter().NumRows()

	// Make the reader using the json string
	jsonStringReader := strings.NewReader(`{"doctorUUID": "00000000-1111-2222-3333-444444444444",
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
	session.Query("DELETE FROM doctors WHERE doctoruuid = 00000000-1111-2222-3333-444444444444").Exec()
}

func TestDoctorGetHandler(t *testing.T) {
	// Doctor info
	doctorUUID, err := gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name := "Test Doctor"
	phone := "123-456-7890"
	primaryFacility := "FakeAddress1"
	primarySpeciality := "Faker1"
	gender := "Male"

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Insert these entries directly into the db
	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
			primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID, name, phone,
		primaryFacility, primarySpeciality, gender).Exec()

	var buff bytes.Buffer
	buff.WriteString("/doctors/doctoruuid/")
	buff.WriteString(doctorUUID.String())
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

	fmt.Println("Deleting Doctor UUID :", doctorUUID.String())
	e := session.Query("DELETE FROM doctors WHERE doctoruuid = ?", doctorUUID).Exec()
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

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
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

func TestPrescriptionCreate(t *testing.T) {
	// Manually add in a fake patient and doctor
	var doctorUUID1 gocql.UUID
	var name1 string
	var phone1 string
	var primaryFacility1 string
	var primarySpeciality1 string
	var gender1 string

	var patientUUID gocql.UUID

	var err error

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Doctor info
	doctorUUID1, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name1 = "Test Doctor"
	phone1 = "123-456-7890"
	primaryFacility1 = "FakeAddress1"
	primarySpeciality1 = "Faker1"
	gender1 = "Male"

	// Patient Info
	patientUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	address := "FakeAddress2"
	bloodType := "O"
	dateOfBirth := "191289601"
	emergencyContact := "415-555-8271"
	gender := "M"
	medicalNumber := "151511517"
	name := "Brown Drey"
	notes := "Broken Legs"
	phone := "151-454-7878"

	// Insert these entries directly into the db
	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
		primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID1, name1, phone1,
		primaryFacility1, primarySpeciality1, gender1).Exec()

	session.Query(`INSERT INTO patients (patientUuid, address, bloodType,
		dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, patientUUID, address, bloodType,
		dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone).Exec()

	numPrescriptionsBefore := session.Query("SELECT * FROM prescriptions").Iter().NumRows()

	var bb bytes.Buffer
	bb.WriteString(`[{"patientUUID":"`)
	bb.WriteString(patientUUID.String())
	bb.WriteString(`","doctorUUID": "`)
	bb.WriteString(doctorUUID1.String())
	bb.WriteString(`","doctor,omitempty": "Dr Ramoray",
			"drug": "Drug Name",
			"startDate": 191289600,
			"endDate": 191389600,
			"instructions,omitempty": "take twice daily"
		}]`)
	entry := bb.String()
	// Make the reader using the json string
	jsonStringReader := strings.NewReader(entry)

	// Create the request with json as body
	req, err := http.NewRequest("POST", "/prescription", jsonStringReader)

	// Check if any errors occured when creating the new request
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(PrescriptionCreate)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	numPrescriptionsAfter := session.Query("SELECT * FROM prescriptions").Iter().NumRows()

	if numPrescriptionsAfter != numPrescriptionsBefore+1 {
		t.Errorf("Expected to have %v prescription, but only have %v", numPrescriptionsBefore+1, numPrescriptionsAfter)
	}

	// CLean up
	e := session.Query("DELETE FROM doctors where doctoruuid = ?", doctorUUID1).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM patients where patientuuid = ?", patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}

}

func TestFutureAppointmentCreateHandler(t *testing.T) {
	var patientUUID gocql.UUID
	var entry string

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = testDB
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
	var patientUUID gocql.UUID
	var appointmentUUID gocql.UUID
	var doctorUUID gocql.UUID
	var err error

	// Doctor info
	doctorUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name := "Test Doctor"
	phone := "123-456-7890"
	primaryFacility := "FakeAddress1"
	primarySpeciality := "Faker1"
	gender := "Male"

	// Patient Info
	patientUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	address := "FakeAddress"
	bloodType := "O"
	dateOfBirth := "191289601"
	emergencyContact := "415-555-8271"
	patientGender := "M"
	medicalNumber := "151511517"
	patientName := "Brown Drey"
	notes := "Broken Legs"
	patientPhone := "151-454-7878"

	// Appointments Info
	appointmentUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	dateScheduled := 1000
	appointmentNotes := "Sample"

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Insert these entries directly into the db
	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
			primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID, name, phone,
		primaryFacility, primarySpeciality, gender).Exec()

	session.Query(`INSERT INTO patients (patientUUID, address, bloodType,
			dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, patientUUID, address, bloodType,
		dateOfBirth, emergencyContact, patientGender, medicalNumber, patientName, notes,
		patientPhone).Exec()

	session.Query(`INSERT INTO futureappointments (appointmentUUID, datescheduled,
		doctoruuid, notes, patientuuid) VALUES (?,?,?,?,?)`,
		appointmentUUID, dateScheduled, doctorUUID, appointmentNotes, patientUUID).Exec()

	var buff bytes.Buffer
	buff.WriteString("/futureappointments/appointmentuuid/")
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
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusOK)
	}

	// Check the body for patientUUID
	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct patientUUID. \nMessage: %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"appointmentUUID":"` + appointmentUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct appointmentUUID. \nMessage: %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"doctorUUID":"` + doctorUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct doctorUUID. \nMessage: %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"notes":"` + appointmentNotes + `"`)) {
		t.Errorf("The response message did not contain the correct notes. \nMessage: %v", rec.Body.String())
	}
	e := session.Query("DELETE FROM patients where patientUUID = ?", patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM doctors where doctorUUID = ?", doctorUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM futureAppointments where appointmentUUID = ?", appointmentUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestCompletedAppointmentCreateHandler(t *testing.T) {
	var patientUUID gocql.UUID

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	// This keyspace can be changed later for tests (i.e. emr_test )
	cluster.Keyspace = testDB
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

	e := session.Query("DELETE FROM completedappointments WHERE appointmentuuid = ?", appointmentUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestCompletedAppointmentGetHandler(t *testing.T) {
	var patientUUID gocql.UUID
	var appointmentUUID gocql.UUID
	var doctorUUID gocql.UUID
	var err error

	// Doctor info
	doctorUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name := "Test Doctor"
	phone := "123-456-7890"
	primaryFacility := "FakeAddress1"
	primarySpeciality := "Faker1"
	gender := "Male"

	// Patient Info
	patientUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	address := "FakeAddress"
	bloodType := "O"
	dateOfBirth := "191289601"
	emergencyContact := "415-555-8271"
	patientGender := "M"
	medicalNumber := "151511517"
	patientName := "Brown Drey"
	notes := "Broken Legs"
	patientPhone := "151-454-7878"

	// Appointments Info
	appointmentUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	dateVisited := 1000
	appointmentNotes := "Sample"
	bloodOxygenLevel := 4
	heartRate := 97
	bloodPressure := 108
	breathingRate := 10

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Insert these entries directly into the db
	e := session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
			primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID, name, phone,
		primaryFacility, primarySpeciality, gender).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query(`INSERT INTO patients (patientUUID, address, bloodType,
			dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, patientUUID, address, bloodType,
		dateOfBirth, emergencyContact, patientGender, medicalNumber, patientName, notes,
		patientPhone).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query(`INSERT INTO completedappointments (appointmentUUID,
		bloodoxygenlevel, bloodpressure, breathingRate, dateVisited, doctoruuid,
		heartrate, notes, patientuuid) VALUES (?,?,?,?,?,?,?,?,?)`, appointmentUUID,
		bloodOxygenLevel, bloodPressure, breathingRate, dateVisited,
		doctorUUID, heartRate, appointmentNotes, patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}

	var buff bytes.Buffer
	buff.WriteString("/completedappointments/appointmentuuid/")
	buff.WriteString(appointmentUUID.String())
	endpoint := buff.String()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RequestURI = endpoint

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(CompletedAppointmentGet)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusOK)
	}

	// Check the body for patientUUID
	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct patientUUID. \nMessage: %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"appointmentUUID":"` + appointmentUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct appointmentUUID. \nMessage: %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"doctorUUID":"` + doctorUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct doctorUUID. \nMessage: %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"notes":"` + appointmentNotes + `"`)) {
		t.Errorf("The response message did not contain the correct notes. \nMessage: %v \n", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"notes":"` + appointmentNotes + `"`)) {
		t.Errorf("The response message did not contain the correct patientUUID. \nMessage: %v", rec.Body.String())
	}
	e = session.Query("DELETE FROM patients where patientUUID = ?", patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM doctors where doctorUUID = ?", doctorUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM completedAppointments where appointmentUUID = ?", appointmentUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestAppointmentGetByDoctorHandler(t *testing.T) {
	var patientUUID gocql.UUID
	var appointmentUUID gocql.UUID
	var appointmentUUID2 gocql.UUID
	var doctorUUID gocql.UUID
	var err error

	// Doctor info
	doctorUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name := "Test Doctor"
	phone := "123-456-7890"
	primaryFacility := "FakeAddress1"
	primarySpeciality := "Faker1"
	gender := "Male"

	// Patient Info
	patientUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	address := "FakeAddress"
	bloodType := "O"
	dateOfBirth := "191289601"
	emergencyContact := "415-555-8271"
	patientGender := "M"
	medicalNumber := "151511517"
	patientName := "Brown Drey"
	notes := "Broken Legs"
	patientPhone := "151-454-7878"

	// Appointments Info
	appointmentUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	dateScheduled := 1000
	appointmentNotes := "Sample"

	// Appointments Info
	appointmentUUID2, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	dateScheduled2 := 2022
	appointmentNotes2 := "Sample2"

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Insert these entries directly into the db
	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
			primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID, name, phone,
		primaryFacility, primarySpeciality, gender).Exec()

	session.Query(`INSERT INTO patients (patientUUID, address, bloodType,
			dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, patientUUID, address, bloodType,
		dateOfBirth, emergencyContact, patientGender, medicalNumber, patientName, notes,
		patientPhone).Exec()

	session.Query(`INSERT INTO futureappointments (appointmentUUID, datescheduled,
		doctoruuid, notes, patientuuid) VALUES (?,?,?,?,?)`,
		appointmentUUID, dateScheduled, doctorUUID, appointmentNotes, patientUUID).Exec()

	session.Query(`INSERT INTO futureappointments (appointmentUUID, datescheduled,
		doctoruuid, notes, patientuuid) VALUES (?,?,?,?,?)`,
		appointmentUUID2, dateScheduled2, doctorUUID, appointmentNotes2, patientUUID).Exec()

	// Get the appointments for patient Brown Drey
	var bb bytes.Buffer
	bb.WriteString("/appointments/doctoruuid/")
	bb.WriteString(doctorUUID.String())

	endpoint := bb.String()
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
	handler := http.HandlerFunc(AppointmentGetByDoctor)
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
	if !strings.Contains(rec.Body.String(), (`"notes":"` + appointmentNotes + `"`)) {
		t.Errorf("The response message did not contain a correct note. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"notes":"` + appointmentNotes2 + `"`)) {
		t.Errorf("The response message did not contain a correct note. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the patientUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"appointmentUUID":"` + appointmentUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the appointmentUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"appointmentUUID":"` + appointmentUUID2.String() + `"`)) {
		t.Errorf("The response message did not contain the appointmentUUID. \n The returned message is: \n %v", rec.Body.String())
	}

	// Clean up the DB
	e := session.Query("DELETE FROM patients where patientuuid = ?", patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM doctors where doctoruuid = ?", doctorUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM futureappointments where appointmentuuid = ?", appointmentUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM futureappointments where appointmentuuid = ?", appointmentUUID2).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestAppointmentGetByPatientHandler(t *testing.T) {
	var patientUUID gocql.UUID
	var appointmentUUID gocql.UUID
	var appointmentUUID2 gocql.UUID
	var doctorUUID gocql.UUID
	var err error

	// Doctor info
	doctorUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name := "Test Doctor"
	phone := "123-456-7890"
	primaryFacility := "FakeAddress1"
	primarySpeciality := "Faker1"
	gender := "Male"

	// Patient Info
	patientUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	address := "FakeAddress"
	bloodType := "O"
	dateOfBirth := "191289601"
	emergencyContact := "415-555-8271"
	patientGender := "M"
	medicalNumber := "151511517"
	patientName := "Brown Drey"
	notes := "Broken Legs"
	patientPhone := "151-454-7878"

	// Appointments Info
	appointmentUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	dateScheduled := 1000
	appointmentNotes := "Sample"

	// Appointments Info
	appointmentUUID2, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	dateScheduled2 := 2022
	appointmentNotes2 := "Sample2"

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Insert these entries directly into the db
	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
			primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID, name, phone,
		primaryFacility, primarySpeciality, gender).Exec()

	session.Query(`INSERT INTO patients (patientUUID, address, bloodType,
			dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, patientUUID, address, bloodType,
		dateOfBirth, emergencyContact, patientGender, medicalNumber, patientName, notes,
		patientPhone).Exec()

	session.Query(`INSERT INTO futureappointments (appointmentUUID, datescheduled,
		doctoruuid, notes, patientuuid) VALUES (?,?,?,?,?)`,
		appointmentUUID, dateScheduled, doctorUUID, appointmentNotes, patientUUID).Exec()

	session.Query(`INSERT INTO futureappointments (appointmentUUID, datescheduled,
		doctoruuid, notes, patientuuid) VALUES (?,?,?,?,?)`,
		appointmentUUID2, dateScheduled2, doctorUUID, appointmentNotes2, patientUUID).Exec()

	// Get the appointments for patient Brown Drey
	var bb bytes.Buffer
	bb.WriteString("/appointments/patientuuid/")
	bb.WriteString(patientUUID.String())

	endpoint := bb.String()
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
	handler := http.HandlerFunc(AppointmentGetByPatient)
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
	if !strings.Contains(rec.Body.String(), (`"notes":"` + appointmentNotes + `"`)) {
		t.Errorf("The response message did not contain a correct note. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"notes":"` + appointmentNotes2 + `"`)) {
		t.Errorf("The response message did not contain a correct note. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the patientUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"appointmentUUID":"` + appointmentUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the appointmentUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"appointmentUUID":"` + appointmentUUID2.String() + `"`)) {
		t.Errorf("The response message did not contain the appointmentUUID. \n The returned message is: \n %v", rec.Body.String())
	}

	// Clean up the DB
	e := session.Query("DELETE FROM patients where patientuuid = ?", patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM doctors where doctoruuid = ?", doctorUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM futureappointments where appointmentuuid = ?", appointmentUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM futureappointments where appointmentuuid = ?", appointmentUUID2).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestPatientGetByDoctorHandler(t *testing.T) {
	var patientUUID gocql.UUID
	var appointmentUUID gocql.UUID
	var doctorUUID gocql.UUID
	var err error

	// Doctor info
	doctorUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name := "Test Doctor"
	phone := "123-456-7890"
	primaryFacility := "FakeAddress1"
	primarySpeciality := "Faker1"
	gender := "Male"

	// Patient Info
	patientUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	address := "FakeAddress"
	bloodType := "O"
	dateOfBirth := "191289601"
	emergencyContact := "415-555-8271"
	patientGender := "M"
	medicalNumber := "151511517"
	patientName := "Brown Drey"
	notes := "Broken Legs"
	patientPhone := "151-454-7878"

	// Appointments Info
	appointmentUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	dateScheduled := 1000
	appointmentNotes := "Sample"

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Insert these entries directly into the db
	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
			primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID, name, phone,
		primaryFacility, primarySpeciality, gender).Exec()

	session.Query(`INSERT INTO patients (patientUUID, address, bloodType,
			dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, patientUUID, address, bloodType,
		dateOfBirth, emergencyContact, patientGender, medicalNumber, patientName, notes,
		patientPhone).Exec()

	session.Query(`INSERT INTO futureappointments (appointmentUUID, datescheduled,
		doctoruuid, notes, patientuuid) VALUES (?,?,?,?,?)`,
		appointmentUUID, dateScheduled, doctorUUID, appointmentNotes, patientUUID).Exec()

	var bb bytes.Buffer
	bb.WriteString("/patients/doctoruuid/")
	bb.WriteString(doctorUUID.String())
	endpoint := bb.String()

	// The doctorUUID is the same as the UUID used for doctors in the test above.
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RequestURI = endpoint

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(PatientGetByDoctor)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusOK)
	}

	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the patientUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"gender":"` + patientGender + `"`)) {
		t.Errorf("The response message did not contain the gender. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"name":"` + patientName + `"`)) {
		t.Errorf("The response message did not contain the name. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"phoneNumber":"` + patientPhone + `"`)) {
		t.Errorf("The response message did not contain the phone number.\n The returned message is: \n %v", rec.Body.String())
	}

	e := session.Query("DELETE FROM patients where patientUUID = ?", patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM doctors where doctorUUID = ?", doctorUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM futureAppointments where appointmentUUID = ?", appointmentUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestDeleteFutureAppointmentHandler(t *testing.T) {
	var patientUUID gocql.UUID
	var appointmentUUID gocql.UUID
	var doctorUUID gocql.UUID

	var e error

	// Doctor info
	doctorUUID, e = gocql.RandomUUID()
	if e != nil {
		t.Fatal(e)
	}
	name := "Test Doctor"
	phone := "123-456-7890"
	primaryFacility := "FakeAddress1"
	primarySpeciality := "Faker1"
	gender := "Male"

	// Patient Info
	patientUUID, e = gocql.RandomUUID()
	if e != nil {
		t.Fatal(e)
	}
	address := "FakeAddress"
	bloodType := "O"
	dateOfBirth := "191289601"
	emergencyContact := "415-555-8271"
	patientGender := "M"
	medicalNumber := "151511517"
	patientName := "Brown Drey"
	notes := "Broken Legs"
	patientPhone := "151-454-7878"

	// Appointments Info
	appointmentUUID, e = gocql.RandomUUID()
	if e != nil {
		t.Fatal(e)
	}
	dateScheduled := 1000
	appointmentNotes := "Sample"

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Insert these entries directly into the db
	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
			primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID, name, phone,
		primaryFacility, primarySpeciality, gender).Exec()

	session.Query(`INSERT INTO patients (patientUuid, address, bloodType,
			dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, patientUUID, address, bloodType,
		dateOfBirth, emergencyContact, patientGender, medicalNumber, patientName, notes,
		patientPhone).Exec()

	session.Query(`INSERT INTO futureappointments (appointmentUUID, datescheduled,
		doctoruuid, notes, patientuuid) VALUES (?,?,?,?,?)`,
		appointmentUUID, dateScheduled, doctorUUID, appointmentNotes, patientUUID).Exec()

	numAppointments := session.Query("SELECT * FROM futureAppointments").Iter().NumRows()

	var bb bytes.Buffer
	bb.WriteString("/futureappointments/appointmentuuid/")
	bb.WriteString(appointmentUUID.String())
	endpoint := bb.String()

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		t.Fatal(err)
	}

	req.RequestURI = endpoint

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(FutureAppointmentDelete)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v but want %v", status, http.StatusOK)
	}

	numAppointments2 := session.Query("SELECT * FROM futureAppointments").Iter().NumRows()

	if numAppointments2+1 != numAppointments {
		t.Errorf("The number of appointments before is %v and the current number of appointments is %v", numAppointments, numAppointments2)
	}

	// Clean up the DB
	e = session.Query("DELETE FROM patients where patientuuid = ?", patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM doctors where doctoruuid = ?", doctorUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM futureappointments where appointmentuuid = ?", appointmentUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestPrescriptionGetByPatient(t *testing.T) {
	var patientUUID gocql.UUID
	var doctorUUID gocql.UUID
	var prescriptionUUID gocql.UUID
	var prescriptionUUID2 gocql.UUID
	var err error

	// Doctor info
	doctorUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	name := "Test Doctor"
	phone := "123-456-7890"
	primaryFacility := "FakeAddress1"
	primarySpeciality := "Faker1"
	gender := "Male"

	// Patient Info
	patientUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	address := "FakeAddress"
	bloodType := "O"
	dateOfBirth := "191289601"
	emergencyContact := "415-555-8271"
	patientGender := "M"
	medicalNumber := "151511517"
	patientName := "Brown Drey"
	notes := "Broken Legs"
	patientPhone := "151-454-7878"

	// Prescription Info
	prescriptionUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	endDate := 191389600

	prescriptionUUID2, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	endDate2 := 200000000

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Insert these entries directly into the db
	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
				primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID, name, phone,
		primaryFacility, primarySpeciality, gender).Exec()

	session.Query(`INSERT INTO patients (patientUuid, address, bloodType,
				dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, patientUUID, address, bloodType,
		dateOfBirth, emergencyContact, patientGender, medicalNumber, patientName, notes,
		patientPhone).Exec()

	err = session.Query(`INSERT INTO prescriptions (patientUUID, prescriptionUUID,
			doctorUUID, doctorName, drug, startDate, endDate, instructions) VALUES (?,?,?,?,?,?,?,?)`,
		patientUUID, prescriptionUUID, doctorUUID, nil, "Wise Drugs", 18019111, endDate, nil).Exec()

	if err != nil {
		t.Fatal(err)
	}

	err = session.Query(`INSERT INTO prescriptions (patientUUID, prescriptionUUID,
				doctorUUID, doctorName, drug, startDate, endDate, instructions) VALUES (?,?,?,?,?,?,?,?)`,
		patientUUID, prescriptionUUID2, doctorUUID, nil, "Bad Drugs", 18011111, endDate2, nil).Exec()

	if err != nil {
		t.Fatal(err)
	}
	// Get the appointments for patient Brown Drey
	var bb bytes.Buffer
	bb.WriteString("/prescriptions/patientuuid/")
	bb.WriteString(patientUUID.String())
	endpoint := bb.String()
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
	handler := http.HandlerFunc(PrescriptionsGetByPatient)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check if the response's uuid is correct (Expected value).
	if !strings.Contains(rec.Body.String(), (`"DoctorUUID":"` + doctorUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the correct doctorUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"PrescriptionUUID":"` + prescriptionUUID.String() + `"`)) {
		t.Errorf("The response message did not contain a correct prescription. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"PrescriptionUUID":"` + prescriptionUUID2.String() + `"`)) {
		t.Errorf("The response message did not contain a correct prescription. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"patientUUID":"` + patientUUID.String() + `"`)) {
		t.Errorf("The response message did not contain the patientUUID. \n The returned message is: \n %v", rec.Body.String())
	}

	// Clean up the DB
	e := session.Query("DELETE FROM patients where patientuuid = ?", patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM doctors where doctoruuid = ?", doctorUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query(`DELETE FROM prescriptions where prescriptionUUID = ? and
			patientUUID = ? and endDate = ?`, prescriptionUUID, patientUUID, endDate).Exec()
	if e != nil {
		t.Fatal(e)
	}

	e = session.Query(`DELETE FROM prescriptions where prescriptionUUID = ? and
			patientUUID = ? and endDate = ?`, prescriptionUUID2, patientUUID, endDate2).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestPatientUpdateHandler(t *testing.T) {
	var patientUUID gocql.UUID
	var err error

	// Patient Info
	patientUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	address := "FakeAddress"
	bloodType := "O"
	dateOfBirth := 191289601
	emergencyContact := "415-555-8271"
	patientGender := "M"
	medicalNumber := "151511517"
	patientName := "Brown Drey"
	notes := "Broken Legs"
	patientPhone := "151-454-7878"

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Add patient to DB
	session.Query(`INSERT INTO patients (patientUuid, address, bloodType,
				dateOfBirth, emergencyContact, gender, medicalNumber, name, notes, phone)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, patientUUID, address, bloodType,
		dateOfBirth, emergencyContact, patientGender, medicalNumber, patientName, notes,
		patientPhone).Exec()

	// Modify a field
	address = "New address"

	// Get the appointments for patient Brown Drey
	var bb bytes.Buffer
	bb.WriteString(`{"patientUUID":"`)
	bb.WriteString(patientUUID.String())
	bb.WriteString(`","address":"`)
	bb.WriteString(address)
	bb.WriteString(`","bloodType":"`)
	bb.WriteString(bloodType)
	bb.WriteString(`","dateOfBirth":`)
	bb.WriteString(strconv.Itoa(dateOfBirth))
	bb.WriteString(`,"emergencyContact":"`)
	bb.WriteString(emergencyContact)
	bb.WriteString(`","gender":"`)
	bb.WriteString(patientGender)
	bb.WriteString(`","medicalNumber":"`)
	bb.WriteString(medicalNumber)
	bb.WriteString(`","name":"`)
	bb.WriteString(patientName)
	bb.WriteString(`","notes":"`)
	bb.WriteString(notes)
	bb.WriteString(`","phoneNumber":"`)
	bb.WriteString(patientPhone)
	bb.WriteString(`"}`)

	entry := bb.String()

	// Make the reader using the json string
	jsonStringReader := strings.NewReader(entry)

	endpoint := "/patients"
	req, err := http.NewRequest("PUT", endpoint, jsonStringReader)
	if err != nil {
		t.Fatal(err)
	}

	req.RequestURI = endpoint
	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(PatientUpdate)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	var patientUUID2 gocql.UUID
	var address2 string
	var bloodType2 string
	var dateOfBirth2 int
	var emergencyContact2 string
	var gender2 string
	var medicalNumber2 string
	var name2 string
	var notes2 string
	var phone2 string

	// Query DB for the patient and check if the changed field (address) is modified
	session.Query("SELECT * FROM patients where patientUUID = ?", patientUUID).Consistency(gocql.One).Scan(&patientUUID2, &address2,
		&bloodType2, &dateOfBirth2, &emergencyContact2, &gender2, &medicalNumber2,
		&name2, &notes2, &phone2)

	// Check fields
	if patientUUID.String() != patientUUID2.String() {
		t.Errorf("PatientUUID did not match. Got %v, expected %v",
			patientUUID2.String(), patientUUID.String())
	}
	if address != address2 {
		t.Errorf("Address did not match. Got %v, expected %v", address2, address)
	}
	if bloodType != bloodType2 {
		t.Errorf("Blood Type did not match. Got %v, expected %v", bloodType2, bloodType)
	}
	if dateOfBirth != dateOfBirth2 {
		t.Errorf("DOB did not match. Got %v, expected %v", dateOfBirth2, dateOfBirth)
	}
	if emergencyContact != emergencyContact2 {
		t.Errorf("Emergency Contact did not match. Got %v, expected %v", emergencyContact2, emergencyContact)
	}
	if patientGender != gender2 {
		t.Errorf("Gender did not match. Got %v, expected %v", gender2, patientGender)
	}
	if medicalNumber != medicalNumber2 {
		t.Errorf("Medical Number did not match. Got %v, expected %v", medicalNumber2, medicalNumber)
	}
	if patientName != name2 {
		t.Errorf("Name did not match. Got %v, expected %v", name2, patientName)
	}
	if notes != notes2 {
		t.Errorf("Notes did not match. Got %v, expected %v", notes2, notes)
	}
	if patientPhone != phone2 {
		t.Errorf("Phone Number did not match. Got %v, expected %v", phone2, patientPhone)
	}

	// Clean up the DB
	e := session.Query("DELETE FROM patients where patientuuid = ?", patientUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestNotificationCreateHandler(t *testing.T) {
	var doctorUUID gocql.UUID
	var doctorUUID2 gocql.UUID
	var err error

	// Doctors info
	doctorUUID, err = gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}
	doctorUUID2, err = gocql.RandomUUID()
	name := "Cyclops"
	name2 := "Wolverine"
	phone := "123-456-7890"
	primaryFacility := "FakeAddress1"
	primarySpeciality := "Faker1"
	gender := "Male"

	// Connect to the database
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
				primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID, name, phone,
		primaryFacility, primarySpeciality, gender).Exec()

	session.Query(`INSERT INTO doctors (doctorUUID, name, phone, primaryFacility,
					primarySpecialty, gender) VALUES (?,?,?,?,?,?)`, doctorUUID2, name2, phone,
		primaryFacility, primarySpeciality, gender).Exec()

	var bb bytes.Buffer
	bb.WriteString(`{"message":"Test Message",`)
	bb.WriteString(`"senderUUID":"`)
	bb.WriteString(doctorUUID2.String())
	bb.WriteString(`","receiverUUID":"`)
	bb.WriteString(doctorUUID.String())
	bb.WriteString(`","senderName":"`)
	bb.WriteString(name)
	bb.WriteString(`"}`)

	entry := bb.String()
	// Make the reader using the json string
	jsonStringReader := strings.NewReader(entry)

	endpoint := "/notifications"
	req, err := http.NewRequest("POST", endpoint, jsonStringReader)
	if err != nil {
		t.Fatal(err)
	}

	req.RequestURI = endpoint
	// Create a response recorder to record the response
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(NotificationCreate)
	handler.ServeHTTP(rec, req)

	status := rec.Code
	if status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusCreated)
	}

	e := session.Query("DELETE FROM doctors WHERE doctorUUID = ?", doctorUUID).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM doctors WHERE doctorUUID = ?", doctorUUID2).Exec()
	if e != nil {
		t.Fatal(e)
	}
}

func TestNotificationsGetByDoctorHandler(t *testing.T) {
	var err error

	// Enter mock data
	dateCreated := 1488254862
	message := "Have you seen Jean?"
	receiverUUID := "4498720b-0491-424f-8e52-6e13bd33da71"
	senderName := "Cyclops"
	senderUUID := "20a5e81c-399f-4777-8bea-9c1fc2388f37"
	notificationUUID, err := gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}

	dateCreated2 := 1388254862
	message2 := "Hey man?"
	notificationUUID2, err := gocql.RandomUUID()
	if err != nil {
		t.Fatal(err)
	}

	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = testDB
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	session.Query(`INSERT INTO notifications (receiverUUID, dateCreated,
		notificationUUID, message, sendername, senderuuid) VALUES (?,?,?,?,?,?)`,
		receiverUUID, dateCreated, notificationUUID, message, senderName,
		senderUUID).Exec()

	session.Query(`INSERT INTO notifications (receiverUUID, dateCreated,
			notificationUUID, message, sendername, senderuuid) VALUES (?,?,?,?,?,?)`,
		receiverUUID, dateCreated2, notificationUUID2, message2, senderName,
		senderUUID).Exec()

	// Get the appointments for patient Brown Drey
	var bb bytes.Buffer
	bb.WriteString("/notifications/doctoruuid/")
	bb.WriteString(receiverUUID)

	endpoint := bb.String()
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
	handler := http.HandlerFunc(NotificationsGetByDoctor)
	handler.ServeHTTP(rec, req)

	// Get the status code of the page and check if it is OK
	status := rec.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	// Check if the response's uuid is correct (Expected value).
	if !strings.Contains(rec.Body.String(), (`"receiverUUID":"` + receiverUUID + `"`)) {
		t.Errorf("The response message did not contain the correct doctorUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"senderUUID":"` + senderUUID + `"`)) {
		t.Errorf("The response message did not contain the correct senderUUID. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"senderName":"` + senderName + `"`)) {
		t.Errorf("The response message did not contain the senderName. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"message":"` + message + `"`)) {
		t.Errorf("The response message did not contain the message. \n The returned message is: \n %v", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), (`"message":"` + message2 + `"`)) {
		t.Errorf("The response message did not contain the message. \n The returned message is: \n %v", rec.Body.String())
	}

	e := session.Query("DELETE FROM notifications where notificationUUID = ? and receiverUUID = ? and datecreated = ?", notificationUUID, receiverUUID, dateCreated).Exec()
	if e != nil {
		t.Fatal(e)
	}
	e = session.Query("DELETE FROM notifications where notificationUUID = ?and receiverUUID = ? and datecreated = ?", notificationUUID2, receiverUUID, dateCreated2).Exec()
	if e != nil {
		t.Fatal(e)
	}

}
