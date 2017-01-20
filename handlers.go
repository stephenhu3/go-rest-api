package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

// Added JSON config file and parser to read from, but removed for demo
const UIAddr = "http://192.168.1.64:3000"
const CASSDB = "127.0.0.1"

func Authenticate(w http.ResponseWriter, r *http.Request) {
	log.Println("OK")
	w.Header().Set("Set-Cookie", "userToken=test; Path=/; HttpOnly")
	w.Header().Set("Access-Control-Allow-Origin", UIAddr)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusOK)
	r.ParseForm()

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", UIAddr)
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var todoId int
	var err error
	if todoId, err = strconv.Atoi(vars["todoId"]); err != nil {
		panic(err)
	}
	todo := RepoFindTodo(todoId)
	if todo.Id > 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound,
		Message: "Not Found"}); err != nil {
		panic(err)
	}

}

/*
Test with this curl command:

curl -H "Content-Type: application/json" -d '{"name":"New Todo"}'
http://localhost:8080/todos

*/
func TodoCreate(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	t := RepoCreateTodo(todo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		panic(err)
	}
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

	log.Printf("Created new patient: %s\t%d\t%s\t%s\t%s",
		patientUUID, gender, medicalNumber, name)

	// insert new patient entry
	if err := session.Query(`INSERT INTO patients (patientUuid, 
		address, bloodType, dateOfBirth, emergencyContact, gender, 
		medicalNumber, name, notes, phone ) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		patientUUID, address, bloodType, dateOfBirth, emergencyContact, 
		gender, medicalNumber, name, notes, phone, ).Exec(); err != nil {
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
		searchUUID).Consistency(gocql.One).Scan( &patientUUID, &address, 
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
		w.WriteHeader(http.StatusFound)
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
Returns a list of patients seen by a specific doctor
Method: GET
Endpoint: //patients/search?doctoruuid=:doctoruuid
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

	var notes string
	var patientUUID gocql.UUID


	// Get all patients that this doctor has seen before
	iter := session.Query("SELECT patientUUID, notes FROM futureappointments WHERE doctoruuid = ?",
		searchUUID).Consistency(gocql.One).Iter();

	// Iterate through returned rows and scan rows 
	// Return unique patients with Unmarshaled JSON string (in Notes) for basic patient info
	m := make(map[gocql.UUID] Patient)
	for iter.Scan( &patientUUID, &notes){
		if _, found := m[patientUUID]; !found{
			currentPatient := Patient{}
			if marshalErr := json.Unmarshal([]byte(notes), &currentPatient); marshalErr != nil{
				log.Println("Patient details malformed at patientuuid", patientUUID)
				continue
			}
			currentPatient.PatientUUID = patientUUID;
			m[patientUUID] = currentPatient
		}
	}

	// Error in iteration returned upon iter.Close()
	// Or no Patients found
	if err := iter.Close(); err !=nil || len(m) == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Status{Code: http.StatusNotFound, Message: "Not Found"})
		log.Printf("PatientLists not found")
		log.Println(err)
		return
	}

	// Patients found
	if len(m) > 0 {
		log.Printf("Patients found")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusFound)

		i := 0
		patientList := make([]Patient, len(m))
		for _, v := range m {
			patientList[i] = v
			i++
		}
		if err := json.NewEncoder(w).Encode(patientList); err != nil {
			panic(err)
		}
	}
}

/*
Create a future appointment
Method: POST
Endpoint: /futureappointment
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
		w.WriteHeader(http.StatusFound)
		if err := json.NewEncoder(w).Encode(FutureAppointment{AppointmentUUID: appointmentUUID,
			PatientUUID: patientUUID, DoctorUUID: doctorUUID, DateScheduled: dateScheduled,
			Notes: notes}); err != nil {
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

	// generate new randomly generated UUID (version 4)
	appointmentUuid, err := gocql.RandomUUID()
	if err != nil {
		log.Fatal(err)
	}
	patientUUID := c.PatientUUID
	doctorUUID := c.DoctorUUID
	dateVisited := c.DateVisited
	breathingRate := c.BreathingRate
	heartRate := c.HeartRate
	bloodOxygenLevel := c.BloodOxygenLevel
	bloodPressure := c.BloodPressure
	notes := c.Notes
	log.Printf("Created completed  appointment: %s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%s\t",
		appointmentUuid, patientUUID, doctorUUID, dateVisited, breathingRate, heartRate,
		bloodOxygenLevel, bloodPressure, notes)

	// insert new completed appointment entry
	if err := session.Query(`INSERT INTO completedAppointments (appointmentUuid,
		patientUUID, doctorUUID, dateVisited, breathingRate, heartRate, bloodOxygenLevel,
		bloodPressure, notes) VALUES (?, ?, ?, ?, ?, ? , ?, ?, ?)`,
		appointmentUuid, patientUUID, doctorUUID, dateVisited, breathingRate, heartRate,
		bloodOxygenLevel, bloodPressure, notes).Exec(); err != nil {
		log.Fatal(err)
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Status{Code: http.StatusCreated,
		Message: "Appointment entry successfully created."})
}

/*
Search for info on a completed appointment
Method: GET
Endpoint: /completedappointments/search?appointmentuuid=:appointmentuuid
*/
func CompletedAppointmentGet(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusFound)
		if err := json.NewEncoder(w).Encode(CompletedAppointment{AppointmentUUID: appointmentUUID,
			PatientUUID: patientUUID, DoctorUUID: doctorUUID, DateVisited: dateVisited,
			BreathingRate: breathingRate, HeartRate: heartRate, BloodOxygenLevel: bloodOxygenLevel,
			BloodPressure: bloodPressure, Notes: notes}); err != nil {
			panic(err)
		}
	}
}
