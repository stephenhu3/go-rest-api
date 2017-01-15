Pull repo to your Go workspace:
```
$GOPATH/src/github.com/{username}
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
-------------------------------------------------------
# API Reference
-------------------------------------------------------

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

HTTP 201 Created

```json
{
  "code": 201,
  "message": "Patient entry successfully created."
}
```
-------------------------------------------------------
GET /patients/search?patientuuid=:patientuuid

**Retrieves a patient record**

Response:

HTTP 302 Found

```json
{
  "patientUUID": "ce3aa844-25cf-4794-9486-83fec2358138",
  "age": 69,
  "gender": "F",
  "insuranceNumber": "Kelly",
  "name": "1234567890"
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

HTTP 201 Created

```json
{
  "code": 201,
  "message": "Appointment entry successfully created."
}
```
-------------------------------------------------------
GET /futureappointments/search?appointmentuuid=:appointmentuuid

**Retrieves a scheduled appointment record**

Response:

HTTP 302 Found

```json
{
  "appointmentUUID": "6b8337bb-b602-4141-aff0-eb52617f1ef9",
  "patientUUID": "36cb5853-758b-44ec-86b4-55cbac3c8afd",
  "doctorUUID": "57c7aea1-9fea-422d-ae35-dbf8ce5f5dda",
  "dateScheduled": 0,
  "notes": "do blood test"
}
```
-------------------------------------------------------
POST {domain}/completedappointments

**Create a completed appointment entry**

Request:

```json
{
    "patientUuid": "36cb5853-758b-44ec-86b4-55cbac3c8afd",
    "doctorUuid": "57c7aea1-9fea-422d-ae35-dbf8ce5f5dda",
    "notes": "do blood test",
	"dateVisited": 1479463552,
	"breathingRate": 10,
	"heartRate": 97,
	"bloodOxygenLevel": 4,
	"bloodPressure": 108
}
```

Response:

HTTP 201 Created

```json
{
  "code": 201,
  "message": "Appointment entry successfully created."
}
```
-------------------------------------------------------
GET /futureappointments/search?appointmentuuid=:appointmentuuid

**Retrieves a completed appointment entry**

Response:

HTTP 302 Found

```json
{
  "appointmentUUID": "4ecafa28-b412-45d9-af2a-758c19bdc433",
  "patientUUID": "36cb5853-758b-44ec-86b4-55cbac3c8afd",
  "doctorUUID": "57c7aea1-9fea-422d-ae35-dbf8ce5f5dda",
  "dateVisited": 1479463552,
  "breathingRate": 10,
  "heartRate": 97,
  "bloodOxygenLevel": 4,
  "bloodPressure": 108,
  "notes": "do blood test"
}
```
-------------------------------------------------------
GET /patients/bydoctor?doctoruuid=:doctoruuid

**Retrieves a list of patients under a doctor with basic info **

Response:

HTTP 200 OK

```json
[
	{
		"patientUuid": "e0736160-82b1-4def-b40b-95f899732024",
		"name": "Kelly",
		"age":"27",
		"gender": "Female",
		"medicalNumber": "1234567890",
		"dateOfBirth": "ISOString",
		"dateOfDeath": "ISOString",
		"ethnicity": "Asian",
		"address": "",
		"phoneNumber": "778-390-4482",
		"notes": "ISOString",
		"emergencyContact": {
			"name": "rick",
			"phoneNumber":"705-190-4482",
			"address": "",
			"relationship": "Brother"
		}
		
	},
	...
]
}
```

-------------------------------------------------------
