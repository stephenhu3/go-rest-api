package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocql/gocql"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

// Added JSON config file and parser to read from, but removed for demo
const UIAddr = "http://192.168.1.64:3000"
const CASSDB = "127.0.0.1"


func PreFlight(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Length", "0")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Range")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("Access-Control-Expose-Headers", "Accept-Ranges, Content-Encoding, Content-Length, Content-Range")

	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode()
}
/*
Validates a user credentials
Method: POST
Endpoint: /login
*/
func UserAuthenticate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	userName := r.Form["username"][0]
	passWord := r.Form["password"][0]

	var userUUID gocql.UUID
	var role string
	var name string

	if err := session.Query(`SELECT userUUID, role, name FROM users WHERE userName = ?
		AND passWord = ?`, userName, passWord).Consistency(gocql.One).Scan(&userUUID,
		&role, &name); err != nil {
		// Incorrect username or password
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Status{Code: http.StatusUnauthorized,
			Message: "Incorrect username or password"})
		log.Printf("Incorrect username or password")
		return
	}

	w.Header().Set("Set-Cookie", "userToken=test; Path=/; HttpOnly")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode( User{ UserUUID: userUUID, Role: role, Name: name } );
	err != nil {
		panic(err)
	}
}

/*
Search for a user's basic info
Method: GET
Endpoint: /users/useruuid/{useruuid}
*/
func UserGet(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	var userUUID gocql.UUID
	var role string
	var name string

	// get the user entry
	if err := session.Query(`SELECT userUUID, role, name FROM users WHERE useruuid=?`,
		searchUUID).Consistency(gocql.One).Scan(&userUUID, &role, &name); err != nil {
		// user not found
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
			Message: "Not Found"})
		log.Printf("User not found")
		log.Println(err)
		return
	}

	// User was found
	if len(userUUID) > 0 {
		log.Printf("User was found")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode( User{ UserUUID: userUUID, Role: role, Name: name } );
		err != nil {
			panic(err)
		}
	}
}

/*
Create a user entry
Method: POST
Endpoint: /users
*/
func UserCreate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	decoder := json.NewDecoder(r.Body)
	var a User
	err := decoder.Decode(&a)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// generate new randomly generated UUID
	userUUID, err := gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}
	userName := a.UserName
	passWord := a.PassWord

	log.Printf("Created new user: %s\t%s\t%s\t",
		userName, passWord, userUUID)

	// insert new user entry
	if err := session.Query(`INSERT INTO users (userName,
		passWord, userUUID) VALUES (?, ?, ?)`, userName, passWord,
		userUUID).Exec(); err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Set-Cookie", "userToken=test; Path=/; HttpOnly")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]gocql.UUID{"userUUID":userUUID})
}

/*
Create a patient entry
Method: POST
Endpoint: /patients
*/
func PatientCreate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	decoder := json.NewDecoder(r.Body)
	var p Patient
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// generate new randomly generated UUID (version 4)
	patientUUID, err := gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}

	address := p.Address
	bloodType := p.BloodType
	dateOfBirth := p.DateOfBirth
	emergencyContact := p.EmergencyContact
	gender := p.Gender
	medicalNumber := p.MedicalNumber
	name := p.Name
	notes := p.Notes
	phone := p.Phone

	log.Printf("Created new patient: %s\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t",
		patientUUID, address, bloodType, dateOfBirth, emergencyContact, gender,
		medicalNumber, name, notes, phone)

	// insert new patient entry
	if err := session.Query(`INSERT INTO patients (patientUuid, 
		address, bloodType, dateOfBirth, emergencyContact, gender, 
		medicalNumber, name, notes, phone ) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		patientUUID, address, bloodType, dateOfBirth, emergencyContact,
		gender, medicalNumber, name, notes, phone).Exec(); err != nil {
		log.Fatal(err)
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Status{Code: http.StatusCreated,
		Message: "Patient entry successfully created."})
}

/*
Search for a patient's info
Method: GET
Endpoint: /patients/search?patientuuid=:patientuuid
*/
func PatientGet(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

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

	// get the patient entry
	if err := session.Query("SELECT * FROM patients WHERE patientUUID = ?",
		searchUUID).Consistency(gocql.One).Scan(&patientUUID, &address,
		&bloodType, &dateOfBirth, &emergencyContact, &gender, &medicalNumber,
		&name, &notes, &phone); err != nil {
		// patient was not found
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
			Message: "Not Found"})
		log.Printf("Patient not found")
		return
	}

	// else, patient was found
	if len(patientUUID) > 0 {
		log.Printf("Patient was found")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Patient{PatientUUID: patientUUID,
			Address: address, BloodType: bloodType, DateOfBirth: dateOfBirth,
			EmergencyContact: emergencyContact, Gender: gender,
			MedicalNumber: medicalNumber, Name: name, Notes: notes,
			Phone: phone}); err != nil {
			panic(err)
		}
	}
}

/*
Update a patient entry
Method: PUT
Endpoint: /patients
*/
func PatientUpdate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	decoder := json.NewDecoder(r.Body)
	var p Patient
	err := decoder.Decode(&p)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	patientUUID := p.PatientUUID
	address := p.Address
	bloodType := p.BloodType
	dateOfBirth := p.DateOfBirth
	emergencyContact := p.EmergencyContact
	gender := p.Gender
	medicalNumber := p.MedicalNumber
	name := p.Name
	notes := p.Notes
	phone := p.Phone

	log.Printf("Udpateing patient: %s\t%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t",
		patientUUID, address, bloodType, dateOfBirth, emergencyContact, gender,
		medicalNumber, name, notes, phone)

	// update patient entry
	if err := session.Query(`UPDATE patients SET address = ?, bloodType = ?, dateOfBirth = ?,
		emergencyContact = ?, gender = ?, medicalNumber = ?, name = ?, notes = ?, phone = ?
		WHERE patientUuid = ? IF EXISTS`,
		address, bloodType, dateOfBirth, emergencyContact, gender, medicalNumber, name, notes,
		phone, patientUUID).Exec(); err != nil {
		// patient was not found
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
			Message: "Error Occured: Patient not updated"})
		log.Printf("Patient not updated")
		return
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Status{Code: http.StatusOK,
		Message: "Patient entry successfully updated."})
}

func mapPatients(m *map[gocql.UUID]string, patientUUID gocql.UUID, session *gocql.Session) {
	// note: need to dereference for map
	if _, found := (*m)[patientUUID]; !found {
		var name string
		// get patient name
		if err := session.Query("SELECT name FROM patients WHERE patientUUID = ?",
			patientUUID).Consistency(gocql.One).Scan(&name); err != nil {
			// patient was not found
			name = "Undefined Patient"
		}
		// cache name to UUID for populating appointment list
		(*m)[patientUUID] = name
	}
}

/*
Returns a list of scheduled and completed appointments for a specific patient
Method: GET
Endpoint: /appointments/patientuuid/{patientuuid}
*/
func AppointmentGetByPatient(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	// Get all future appointments by patient
	iter := session.Query("SELECT * FROM futureappointments WHERE patientuuid = ?",
		searchUUID).Consistency(gocql.One).Iter()

	// Get all completed appointments by patient
	completedIter := session.Query("SELECT * FROM completedappointments WHERE patientuuid = ?",
		searchUUID).Consistency(gocql.One).Iter()

	// no appointments found
	if iter.NumRows() == 0 && completedIter.NumRows() == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound, Message: "Not Found"})
		log.Printf("Scheduled appointments not found")
		return
	}

	appointmentList := make([]GenericAppointment, iter.NumRows()+completedIter.NumRows())
	i := 0
	m := make(map[gocql.UUID]string)
	var appointmentUUID gocql.UUID
	var dateVisited int
	var dateScheduled int
	var doctorUUID gocql.UUID
	var notes string
	var patientUUID gocql.UUID
	var name string

	// scheduled appointment(s) found
	if iter.NumRows() > 0 {
		log.Printf("Scheduled appointments found")

		for iter.Scan(&appointmentUUID, &dateScheduled, &doctorUUID, &notes, &patientUUID) {
			// Search patient table to get patient name, cache patient names
			mapPatients(&m, patientUUID, session)
			name = m[patientUUID]

			appointmentList[i] = GenericAppointment{
				AppointmentUUID: appointmentUUID, PatientUUID: patientUUID,
				DoctorUUID: doctorUUID, DateScheduled: dateScheduled,
				DateVisited: 0, Notes: notes, PatientName: name}
			i++
		}
	}

	// completed appointment(s) found
	if completedIter.NumRows() > 0 {
		log.Printf("Completed appointments found")

		for completedIter.Scan(&appointmentUUID, nil, nil, nil, &dateVisited,
			&doctorUUID, nil, &notes, &patientUUID) {
			mapPatients(&m, patientUUID, session)
			name = m[patientUUID]

			appointmentList[i] = GenericAppointment{
				AppointmentUUID: appointmentUUID, PatientUUID: patientUUID,
				DoctorUUID: doctorUUID, DateScheduled: 0, DateVisited: dateVisited,
				Notes: notes, PatientName: name}
			i++
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(appointmentList); err != nil {
		panic(err)
	}
}

/*
Returns a list of scheduled and completed appointments for a specific doctor
Method: GET
Endpoint: /appointments/doctoruuid/{doctoruuid}
*/
func AppointmentGetByDoctor(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	// Get all future appointments by doctor
	iter := session.Query("SELECT * FROM futureappointments WHERE doctoruuid = ?",
		searchUUID).Consistency(gocql.One).Iter()

	// Get all completed appointments by doctor
	completedIter := session.Query("SELECT * FROM completedappointments WHERE doctoruuid = ?",
		searchUUID).Consistency(gocql.One).Iter()

	// no appointments found
	if iter.NumRows() == 0 && completedIter.NumRows() == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound, Message: "Not Found"})
		log.Printf("Scheduled appointments not found")
		return
	}

	appointmentList := make([]GenericAppointment, iter.NumRows()+completedIter.NumRows())
	i := 0
	m := make(map[gocql.UUID]string)
	var appointmentUUID gocql.UUID
	var dateVisited int
	var dateScheduled int
	var doctorUUID gocql.UUID
	var notes string
	var patientUUID gocql.UUID
	var name string

	// scheduled appointment(s) found
	if iter.NumRows() > 0 {
		log.Printf("Scheduled appointments found")

		for iter.Scan(&appointmentUUID, &dateScheduled, &doctorUUID, &notes, &patientUUID) {
			// Search patient table to get patient name, cache patient names
			// TODO Optimization: create table of appointments by doctor
			mapPatients(&m, patientUUID, session)
			name = m[patientUUID]

			appointmentList[i] = GenericAppointment{
				AppointmentUUID: appointmentUUID, PatientUUID: patientUUID,
				DoctorUUID: doctorUUID, DateScheduled: dateScheduled,
				DateVisited: 0, Notes: notes, PatientName: name}
			i++
		}
	}

	// completed appointment(s) found
	if completedIter.NumRows() > 0 {
		log.Printf("Completed appointments found")

		for completedIter.Scan(&appointmentUUID, nil, nil, nil, &dateVisited,
			&doctorUUID, nil, &notes, &patientUUID) {
			mapPatients(&m, patientUUID, session)
			name = m[patientUUID]

			appointmentList[i] = GenericAppointment{
				AppointmentUUID: appointmentUUID, PatientUUID: patientUUID,
				DoctorUUID: doctorUUID, DateScheduled: 0, DateVisited: dateVisited,
				Notes: notes, PatientName: name}
			i++
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(appointmentList); err != nil {
		panic(err)
	}
}

/*
Returns a list of patients seen by a specific doctor
Method: GET
Endpoint: /patients/doctoruuid/{doctoruuid}
*/
func PatientGetByDoctor(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	// Get all future appointments by doctor
	iter := session.Query("SELECT * FROM futureappointments WHERE doctoruuid = ?",
		searchUUID).Consistency(gocql.One).Iter()

	// Get all completed appointments by doctor
	completedIter := session.Query("SELECT * FROM completedappointments WHERE doctoruuid = ?",
		searchUUID).Consistency(gocql.One).Iter()

	// no appointments found, thus no patients
	if iter.NumRows() == 0 && completedIter.NumRows() == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound, Message: "Not Found"})
		log.Printf("Patients by doctor not found")
		return
	}

	// make a set of patientUUIDs
	m := make(map[gocql.UUID]gocql.UUID)
	var patientUUID gocql.UUID

	// scheduled appointment(s) found
	if iter.NumRows() > 0 {
		log.Printf("Scheduled appointments found")
		for iter.Scan(nil, nil, nil, nil, &patientUUID) {
			m[patientUUID] = patientUUID
		}
	}

	// completed appointment(s) found
	if completedIter.NumRows() > 0 {
		log.Printf("Completed appointments found")
		for completedIter.Scan(nil, nil, nil, nil, nil, nil, nil, nil, &patientUUID) {
			m[patientUUID] = patientUUID
		}
	}

	var patientList []Patient
	var dateOfBirth int
	var gender string
	var name string
	var phone string

	// get each patient's info and add to list
	for k := range m {
		if err := session.Query("SELECT * FROM patients WHERE patientUUID = ?",
			k).Consistency(gocql.One).Scan(&patientUUID, nil,
			nil, &dateOfBirth, nil, &gender, nil, &name, nil, &phone); err != nil {
			log.Printf("Patient does not exist, skipping")
		} else {
			patientList = append(patientList, Patient{PatientUUID: patientUUID,
				DateOfBirth: dateOfBirth, Gender: gender, Name: name, Phone: phone})
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(patientList); err != nil {
		panic(err)
	}
}

/*
Create a future appointment
Method: POST
Endpoint: /futureappointments
*/
func FutureAppointmentCreate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	decoder := json.NewDecoder(r.Body)
	var f FutureAppointment
	err := decoder.Decode(&f)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// generate new randomly generated UUID (version 4)
	appointmentUuid, err := gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}
	patientUUID := f.PatientUUID
	doctorUUID := f.DoctorUUID
	dateScheduled := f.DateScheduled
	notes := f.Notes
	log.Printf("Created future appointment: %s\t%d\t%s\t%s\t%s",
		appointmentUuid, patientUUID, doctorUUID, dateScheduled, notes)

	// insert new appointment entry
	if err := session.Query(`INSERT INTO futureAppointments (appointmentUuid,
		patientUUID, doctorUUID, dateScheduled, notes) VALUES (?, ?, ?, ?, ?)`,
		appointmentUuid, patientUUID, doctorUUID, dateScheduled, notes).Exec(); err != nil {
		log.Fatal(err)
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Status{Code: http.StatusCreated,
		Message: "Appointment entry successfully created."})
}

/*
Search for info on a future appointment
Method: GET
Endpoint: /futureappointments/search?appointmentuuid=:appointmentuuid
*/
func FutureAppointmentGet(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	var searchUUID = r.Form["appointmentuuid"][0]

	var appointmentUUID gocql.UUID
	var patientUUID gocql.UUID
	var doctorUUID gocql.UUID
	var dateScheduled int
	var notes string

	// get the appointment entry
	if err := session.Query("SELECT * FROM futureAppointments WHERE appointmentUUID = ?",
		searchUUID).Consistency(gocql.One).Scan(&appointmentUUID, &dateScheduled,
		&doctorUUID, &notes, &patientUUID); err != nil {
		// appointment was not found
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
			Message: "Not Found"})
		log.Printf("Appointment not found")
		return
	}

	// else, appointment was found
	if len(appointmentUUID) > 0 {
		log.Printf("Appointment was found")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(FutureAppointment{
			AppointmentUUID: appointmentUUID, PatientUUID: patientUUID,
			DoctorUUID: doctorUUID, DateScheduled: dateScheduled, Notes: notes}); err != nil {
			panic(err)
		}
	}
}

/*
Create a completed appointment
Method: POST
Endpoint: /completedappointments
*/
func CompletedAppointmentCreate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	decoder := json.NewDecoder(r.Body)
	var c CompletedAppointment
	err := decoder.Decode(&c)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	appointmentUuid := c.AppointmentUUID
	patientUUID := c.PatientUUID
	doctorUUID := c.DoctorUUID
	dateVisited := c.DateVisited
	breathingRate := c.BreathingRate
	heartRate := c.HeartRate
	bloodOxygenLevel := c.BloodOxygenLevel
	bloodPressure := c.BloodPressure
	notes := c.Notes

	// var deleteSuccess bool
	if _, err := session.Query(`DELETE FROM futureappointments WHERE appointmentuuid=? IF EXISTS`,
		appointmentUuid).ScanCAS(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Updating appointment: %s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%s\t",
		appointmentUuid, patientUUID, doctorUUID, dateVisited, breathingRate, heartRate,
		bloodOxygenLevel, bloodPressure, notes)

	// update appointment entry, create entry if does not exist
	if err := session.Query(`UPDATE completedappointments SET patientUUID = ?,
		doctorUUID = ?, dateVisited = ?, breathingRate = ?, heartRate = ?, bloodOxygenLevel = ?,
		bloodPressure = ?, notes = ? WHERE appointmentUuid = ?`, patientUUID,
		doctorUUID, dateVisited, breathingRate, heartRate, bloodOxygenLevel, bloodPressure, notes,
		appointmentUuid).Exec(); err != nil {
		// Appointment not created/updated
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
			Message: "Error Occured: Appointment not updated/created"})
		log.Printf("Appointment not updated/created")
		log.Println(err)
		return
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Status{Code: http.StatusCreated,
		Message: "Appointment entry successfully updated/created."})
}

/*
Search for info on a completed appointment
Method: GET
Endpoint: /completedappointments/appointmentuuid/{appointmentuuid}
*/
func CompletedAppointmentGet(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()


	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	var appointmentUUID gocql.UUID
	var patientUUID gocql.UUID
	var doctorUUID gocql.UUID
	var dateVisited int
	var breathingRate int
	var heartRate int
	var bloodOxygenLevel int
	var bloodPressure int
	var notes string

	// get the appointment entry (match arguments with alphabetical positioning of retrieved columns)
	if err := session.Query("SELECT * FROM completedAppointments WHERE appointmentUUID = ?",
		searchUUID).Consistency(gocql.One).Scan(&appointmentUUID, &bloodOxygenLevel, &bloodPressure,
		&breathingRate, &dateVisited, &doctorUUID, &heartRate, &notes, &patientUUID); err != nil {
		// appointment was not found
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
			Message: "Not Found"})
		log.Printf("Appointment not found")
		return
	}

	// else, appointment was found
	if len(appointmentUUID) > 0 {
		log.Printf("Appointment was found")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(CompletedAppointment{AppointmentUUID: appointmentUUID,
			PatientUUID: patientUUID, DoctorUUID: doctorUUID, DateVisited: dateVisited,
			BreathingRate: breathingRate, HeartRate: heartRate, BloodOxygenLevel: bloodOxygenLevel,
			BloodPressure: bloodPressure, Notes: notes}); err != nil {
			panic(err)
		}
	}
}

/*
Delete selected appointment
Method: DELETE
Endpoint: /futureappointments/appointmentuuid/{appointmentuuid}
*/
func FutureAppointmentDelete(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()


	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	// Tries to delete from futureAppointments
	if deleteSuccess, err := session.Query("DELETE FROM futureAppointments WHERE appointmentuuid=? IF EXISTS",
		searchUUID).ScanCAS(); err != nil || !deleteSuccess {
		log.Println(deleteSuccess)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
				Message: "Delete target not found"}); err != nil {
			panic(err)
		}
	} else {
		log.Println(deleteSuccess)
		log.Printf("Delete on: %s\t", searchUUID)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Status{Code: http.StatusOK,
				Message: "Delete Success"}); err != nil {
			panic(err)
		}
	}
}

/*
Create a Doctor entry
Method: POST
Endpoint: /doctors
*/
func DoctorCreate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	decoder := json.NewDecoder(r.Body)
	var d Doctor
	err := decoder.Decode(&d)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	doctorUUID := d.DoctorUUID
	name := d.Name
	phoneNumber := d.Phone
	prFacility := d.PrimaryFacility
	prSpecialty := d.PrimarySpecialty
	gender := d.Gender

	log.Printf("Created new doctor: %s\t%s\t%s\t%s\t%s\t%s\t",
		doctorUUID, name, phoneNumber, prFacility, prSpecialty, gender)

	// insert new doctor entry
	if err := session.Query(`INSERT INTO doctors (doctorUUID,
		name, phone, primaryFacility, primarySpecialty,
		gender) VALUES (?, ?, ?, ?, ?, ?)`,
		doctorUUID, name, phoneNumber, prFacility, prSpecialty,
		gender).Exec(); err != nil {
		log.Fatal(err)
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Status{Code: http.StatusCreated,
		Message: "Doctor entry successfully created."})
}

/*
Search for a doctor's info
Method: GET
Endpoint: /doctors/doctoruuid/{doctoruuid}
*/
func DoctorGet(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	var doctorUUID gocql.UUID
	var name string
	var phoneNumber string
	var primaryFacility string
	var primarySpecialty string
	var gender string

	// get the doctor entry
	if err := session.Query("SELECT * FROM doctors WHERE doctorUUID = ?",
		searchUUID).Consistency(gocql.One).Scan(&doctorUUID, &gender,
		&name, &phoneNumber, &primaryFacility, &primarySpecialty); err != nil {
		// doctor was not found
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
			Message: "Not Found"})
		log.Printf("No Doctor found")
		return
	}

	// else, doctor was found
	if len(doctorUUID) > 0 {
		log.Printf("Doctor was found")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(Doctor{DoctorUUID: doctorUUID,
			Name: name, Phone: phoneNumber, PrimaryFacility: primaryFacility,
			PrimarySpecialty: primarySpecialty, Gender: gender}); err != nil {
			panic(err)
		}
	}
}

/*
Returns a list of all doctors
Method: GET
Endpoint: /doctors
*/
func DoctorListGet(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// Get all doctors of current clinic
	iter := session.Query("SELECT * FROM doctors").Consistency(gocql.One).Iter()

	doctorList := make([]Doctor, iter.NumRows())
	i := 0
	var doctorUUID gocql.UUID
	var name string
	var phoneNumber string
	var primaryFacility string
	var primarySpecialty string
	var gender string

	// doctors found
	if iter.NumRows() > 0 {
		log.Printf("Scheduled appointments found")

		for iter.Scan(&doctorUUID, &gender, &name, &phoneNumber, &primaryFacility,
			&primarySpecialty) {
			doctorList[i] = Doctor{DoctorUUID: doctorUUID,
			Name: name, Phone: phoneNumber, PrimaryFacility: primaryFacility,
			PrimarySpecialty: primarySpecialty, Gender: gender}
			i++
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(doctorList); err != nil {
		panic(err)
	}
}

/*
Returns a list of prescriptions for a specific patient
Method: GET
Endpoint: /prescriptions/patientuuid/{patientuuid}
*/
func PrescriptionsGetByPatient(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	// Get all prescriptions for a patient
	iter := session.Query(`SELECT * FROM prescriptions WHERE patientuuid = ?
		ORDER BY endDate DESC`, searchUUID).Consistency(gocql.One).Iter()

	prescriptionList := make(Prescriptions, iter.NumRows())
	i := 0

	var doctorName string
	var doctorUUID gocql.UUID
	var drug string
	var endDate int
	var instructions string
	var patientUUID gocql.UUID
	var prescriptionUUID gocql.UUID
	var startDate int

	// prescriptions found
	if iter.NumRows() > 0 {
		log.Printf("prescriptions found")

		for iter.Scan(&patientUUID, &endDate, &prescriptionUUID, &doctorName,
			&doctorUUID, &drug, &instructions, &startDate){
			prescriptionList[i] = Prescription{
				DoctorName: doctorName, DoctorUUID: doctorUUID, Drug: drug,
				EndDate: endDate, Instructions: instructions, PatientUUID: patientUUID,
				PrescriptionUUID: prescriptionUUID, StartDate: startDate}
			i++
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(prescriptionList); err != nil {
		panic(err)
	}
}

/*
Create a new Prescription
Method: POST
Endpoint: /prescription
*/
func PrescriptionCreate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	decoder := json.NewDecoder(r.Body)
	var prescriptionList Prescriptions
	err := decoder.Decode(&prescriptionList)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	for _, d := range prescriptionList {
		// generate new randomly generated UUID
		prescriptionUUID, err := gocql.RandomUUID()
		if err != nil {
			log.Fatal(err)
		}

		doctorName := d.DoctorName
		doctorUUID := d.DoctorUUID
		drug := d.Drug
		endDate := d.EndDate
		instructions := d.Instructions
		patientUUID := d.PatientUUID
		startDate := d.StartDate

		log.Printf("Created new prescription: %s\t%s\t%s\t%d\t%s\t%s\t%s\t%d\t",
			doctorName, doctorUUID, drug, endDate, instructions, patientUUID,
			prescriptionUUID, startDate)

		// insert new prescription entry
		if err := session.Query(`INSERT INTO prescriptions (doctorName, doctorUUID,
			drug, endDate, instructions, patientUUID, prescriptionUUID, startDate)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, doctorName, doctorUUID, drug, endDate,
			instructions, patientUUID, prescriptionUUID, startDate).Exec(); err != nil {
			log.Fatal(err)
		}
	}
	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Status{Code: http.StatusCreated,
		Message: "Prescription entry successfully created."})
}

/*
Create a new notification for a doctor
Method: POST
Endpoint: /notifications
*/
func NotificationCreate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	decoder := json.NewDecoder(r.Body)
	var n Notification
	err := decoder.Decode(&n)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// generate new randomly generated UUID
	notificationUUID, err := gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}

	date := n.Date
	message := n.Messsage
	receiverUUID := n.ReceiverUUID
	senderName := n.SenderName
	senderUUID := n.SenderUUID

	log.Printf("Creating new notification: %d\t%s\t%s\t%s\t%s\t%s\t",
		date, message, senderUUID, receiverUUID, senderName, senderUUID)

	// insert new notification entry
	if err := session.Query(`INSERT INTO notifications (receiverUUID, date, notificationUUID,
		message, senderName, senderUUID) VALUES (?, ?, ?, ?, ?, ?)`,
		receiverUUID, date, notificationUUID, message, senderName, senderUUID).Exec();
		err != nil {
			log.Fatal(err)
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Status{Code: http.StatusCreated,
		Message: "Notification entry successfully created."})
}


/*
Returns a list of notifications for a specific doctor
Method: GET
Endpoint: /notifications/doctoruuid/{doctoruuid}
*/
func NotificationsGetByDoctor(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	// Get all notifications for a doctor, limit to last 100
	iter := session.Query(`SELECT receiverUUID, date, message, senderUUID, senderName FROM
		notifications WHERE receiverUUID = ? LIMIT 100`, searchUUID).Consistency(gocql.One).Iter()

	notiList := make(Notifications, iter.NumRows())
	i := 0

	var date int
	var message string
	var receiverUUID gocql.UUID
	var senderName string
	var senderUUID gocql.UUID

	// notifications found
	if iter.NumRows() > 0 {
		log.Printf("Notifications found")

		for iter.Scan(&receiverUUID, &date, &message, &senderUUID, &senderName){
			notiList[i] = Notification {
				Date: date, Messsage: message,
				ReceiverUUID: receiverUUID, SenderName: senderName, SenderUUID: senderUUID }
			i++
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(notiList); err != nil {
		panic(err)
	}
}

/*

Upload a document
Method: POST
Endpoint: /document
*/
func DocumentCreate(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	// generate new randomly generated UUID (version 4)
	documentUUID, err := gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}

	patientUUID := r.FormValue("patientUUID")
	filename := r.FormValue("filename")
	dateUploaded := r.FormValue("dateUploaded")

	file, _, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	binaryContent, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal("error:", err)
	}

	log.Printf("Created new document: %s\t%s\t%s\t%d\t",
		documentUUID, patientUUID, filename, dateUploaded)

	// insert new document entry
	if err := session.Query(`INSERT INTO documents (documentUUID,
		patientUUID, filename, dateUploaded, content) VALUES (?, ?, ?, ?, ?)`,
		documentUUID, patientUUID, filename, dateUploaded,
		binaryContent).Exec(); err != nil {
		log.Fatal(err)
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	// returns the created documentuuid on success
	json.NewEncoder(w).Encode(map[string]gocql.UUID{"documentuuid":documentUUID})
}

/*
Download a document given its UUID
Method: GET
Endpoint: /documents/documentuuid/{documentuuid}
*/
func DocumentGet(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	var documentUUID gocql.UUID
	var filename string
	var content []byte

	// download the document
	if err := session.Query("SELECT content, filename FROM documents WHERE documentUUID = ?",
		searchUUID).Consistency(gocql.One).Scan(&content, &filename); err != nil {
		// document was not found
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
			Message: "Not Found"})
		log.Printf("No document found")
		return
	}

	// else, document was found
	if len(documentUUID) > 0 {
		log.Printf("document was found")
		err := ioutil.WriteFile(filename, content, 0755)
		// issue with writing file
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
				Message: "Not Found"})
			return
		}

		fopen, err := os.Open(filename)
		defer fopen.Close()
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
				Message: "Not Found"})
			return
		}

		// determine file content type
		header := make([]byte, 512)
		fopen.Read(header)
		fileType := http.DetectContentType(header)

		stat, _ := fopen.Stat()
		fileSize := strconv.FormatInt(stat.Size(), 10)

		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", fileType)
		w.Header().Set("Content-Length", fileSize)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Range")
		w.Header().Set("Access-Control-Expose-Headers", "Accept-Ranges, Content-Encoding, Content-Length, Content-Range")

		w.WriteHeader(http.StatusOK)

		// send the file (read 512 bytes from the file already so reset offset)
		fopen.Seek(0, 0)
		io.Copy(w, fopen)
		return
	}
}

/*
Returns an index list of documents for a given patient
Method: GET
Endpoint: /documents/patientuuid/{patientuuid}
*/
func DocumentListGetByPatient(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster(CASSDB)
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	if URI := strings.Split(r.RequestURI, "/"); len(URI) != 4 {
		panic("Improper URI")
	}

	var searchUUID = strings.Split(r.RequestURI, "/")[3]

	// Get all documents metadata of a patient
	iter := session.Query(`SELECT documentuuid, dateuploaded, filename, patientuuid FROM documents
		WHERE patientuuid = ?`, searchUUID).Consistency(gocql.One).Iter()

	// no documents found
	if iter.NumRows() == 0 {
		log.Printf("No documents found for patient")
	}

	docuList := make([]Document, iter.NumRows())
	i := 0
	var documentUUID gocql.UUID
	var patientUUID gocql.UUID
	var filename     string
	var dateUploaded int

	// documents found
	if iter.NumRows() > 0 {
		log.Printf("Documents found")

		for iter.Scan(&documentUUID, &dateUploaded, &filename, &patientUUID) {
			docuList[i] = Document{
				 DocumentUUID: documentUUID, PatientUUID: patientUUID, Filename: filename,
				 DateUploaded: dateUploaded}
			i++
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(docuList); err != nil {
		panic(err)
	}
}