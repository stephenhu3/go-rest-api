Pull repo to your Go workspace:
```
$GOPATH/src/github.com/{username}
```

Requirement:
Gorilla MUX
```
go get github.com/gorilla/mux
```

Install the program to the bin directory of your workspace:
```
cd $GOPATH/src/
go install github.com/{username}/go-rest
```


Run the program:
```
$GOPATH/bin/go-rest
```

API Reference

POST {domain}/patients

**Create a new patient**

Request:

```json
{
	"age": 69,
	"gender": "F",
	"name": "Kelly",
	"insuranceNumber": "1234567890"
}
```

Response:

HTTP 200

```json
{
	"status": "success"
}
```

-------------------------------------------------------

POST {domain}/futureappointments

**Create a future appointment**

Request:

```json
{
	"patientUuid": "219529de-0c17-431b-8363-6fcb32e2f708",
	"doctorUuid": "57c7aea1-9fea-422d-ae35-dbf8ce5f5dda",
	"dateScheduled": 1479463552,
	"notes": "do blood test"
}
```

Response:

HTTP 200

```json
{
	"status": "success"
}
```

-------------------------------------------------------

GET /patients/search?patientuuid=:patientuuid

**Retrieves a patient record**

Response:

HTTP 200

```json

	"patientUuid": "e0736160-82b1-4def-b40b-95f899732024",
	"age": "27",
	"gender": "Female",
	"insuranceNumber": "1234567890",
	"name": "Kelly"
}
```

-------------------------------------------------------
