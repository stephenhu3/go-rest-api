package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound,
		Text: "Not Found"}); err != nil {
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
	cluster := gocql.NewCluster("127.0.0.1")
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
	age := p.Age
	gender := p.Gender
	insuranceNumber := p.InsuranceNumber
	name := p.Name
	log.Printf("Created new patient: %s\t%d\t%s\t%s\t%s",
		patientUUID, age, gender, insuranceNumber, name)

	// insert new patient entry
	if err := session.Query(`INSERT INTO patients (patientUuid, age, gender,
		name, insuranceNumber) VALUES (?, ?, ?, ?, ?)`,
		patientUUID, age, gender, insuranceNumber, name).Exec(); err != nil {
		log.Fatal(err)
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("{\"status\": \"success\"}")
}

/*
Search for a patient's info
Method: POST
Endpoint: /patients/search
*/
func PatientSearch(w http.ResponseWriter, r *http.Request) {
	// connect to the cluster
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "emr"
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	defer session.Close()

	err := r.ParseForm()
    if err != nil {
       panic(err)
    }
    var searchUUID = r.Form["PatientUUID"][0]

	var patientUUID gocql.UUID
	var age int
	var gender string
	var insuranceNumber string
	var name string

	// get the patient entry
	if err := session.Query("SELECT * FROM patients WHERE patientUUID = ?",
		searchUUID).Consistency(gocql.One).Scan(&patientUUID, &age, &gender,
		&insuranceNumber, &name); err != nil {
		// patient was not found
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound,
			Text: "Not Found"})
		// w.WriteHeader(http.StatusNotFound)
		log.Printf("Patient not found")
		return
	}

	// else, patient was found
	// TODO: Fix nil fields in response
	if len(patientUUID) > 0 {
		log.Printf("Patient was found")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := json.NewEncoder(w).Encode(Patient{PatientUUID: patientUUID,
			Age: age, Gender: gender, InsuranceNumber: insuranceNumber,
			Name: name}); err != nil {
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
	cluster := gocql.NewCluster("127.0.0.1")
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
	patientUuid := f.PatientUuid
	doctorUuid := f.DoctorUuid
	dateScheduled := f.DateScheduled
	notes := f.Notes
	log.Printf("Created future appointment: %s\t%d\t%s\t%s\t%s",
		appointmentUuid, patientUuid, doctorUuid, dateScheduled, notes)

	// insert new appointment entry
	if err := session.Query(`INSERT INTO futureAppointments (appointmentUuid,
		patientUuid, doctorUuid, dateScheduled, notes) VALUES (?, ?, ?, ?, ?)`,
		appointmentUuid, patientUuid, doctorUuid, dateScheduled, notes).Exec(); err != nil {
		log.Fatal(err)
	}

	// send success response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("{\"status\": \"success\"}")
}
