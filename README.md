Pull repo to your Go workspace:
```
$GOPATH/src/github.com/{username}
```

Install the program to the bin directory of your workspace:
```
cd $GOPATH/src/
go install github.com/{username}/go-rest
```

Set up Database:
```
cqlsh --request-timeout 120 localhost < cqlsh-setup.cql
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
  "address": "5698 Cedar Avenue, San Francisco, California",
  "bloodType": "B-Positive",
  "dateOfBirth": 191289600,
  "emergencyContact": "415-555-8271",
  "gender": "F",
  "medicalNumber": "1234567890",
  "name": "Kelly Lai",
  "notes": "Accompanied by guide dog, ensure patient's area is wheelchair-friendly",
  "phoneNumber": "483-555-5123"
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
GET /patients/patientuuid/{patientuuid}

**Retrieves a patient record**

Response:

HTTP 302 Found

```json
{
  "address": "5698 Cedar Avenue, San Francisco, California",
  "bloodType": "B-Positive",
  "dateOfBirth": 191289600,
  "emergencyContact": "415-555-8271",
  "gender": "F",
  "medicalNumber": "1234567890",
  "name": "Kelly Lai",
  "notes": "Accompanied by guide dog, ensure patient's area is wheelchair-friendly",
  "phoneNumber": "483-555-5123"
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
GET /futureappointments/search?appointmentuuid={appointmentuuid}

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
DELETE /futureappointments/appointmentuuid/{appointmentuuid}

**Deletes a futureappoitnment entry**

Response:

HTTP 200 Found

```json
{
  "code": 200,
  "message": "Delete Success"
}
```

```json

HTTP 401 NotFound
{
  "code": 401,
  "message": "Delete target not found"
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
GET /completedappointments/search?appointmentuuid=:appointmentuuid

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

GET /patients/doctoruuid/{doctoruuid}

**Retrieves a list of patients and their basic info that have been treated by, or is scheduled with the doctor**

Response:

HTTP 302 Found

```json
[
  {
    "patientUUID": "e572fe98-4662-47f7-930c-cf4f7d13e26e",
    "dateOfBirth": 191289600,
    "gender": "M",
    "name": "Lamar Odom",
    "phoneNumber": "483-555-5123"
  },
  {
    "patientUUID": "6c202b83-1759-40e2-b1c8-7566330366d2",
    "dateOfBirth": 191289600,
    "gender": "F",
    "name": "Kelly Lai",
    "phoneNumber": "483-555-5123"
  },
  {
    "patientUUID": "6e894f6b-cbf6-4703-ad4f-bd93126450cb",
    "dateOfBirth": 191289600,
    "gender": "M",
    "name": "Joey Kapow",
    "phoneNumber": "483-555-5123"
  }
]
```

-------------------------------------------------------

GET /appointments/doctoruuid/{doctoruuid}

**Retrieves a list of scheduled and completed appointments under a doctor**

Response:

HTTP 302 Found

```json
[
  {
    "appointmentUUID": "30c40285-bbf5-4a09-b849-b9aa4c0f9f97",
    "patientUUID": "e572fe98-4662-47f7-930c-cf4f7d13e26e",
    "doctorUUID": "57c7aea1-9fea-422d-ae35-dbf8ce5f5dda",
    "patientName": "Lamar Odom",
    "dateScheduled": 1479463552,
    "dateVisited": 0,
    "notes": "do something"
  },
  {
    "appointmentUUID": "13c57da1-6f13-4898-b42b-de4252131337",
    "patientUUID": "6c202b83-1759-40e2-b1c8-7566330366d2",
    "doctorUUID": "57c7aea1-9fea-422d-ae35-dbf8ce5f5dda",
    "patientName": "Kelly Lai",
    "dateScheduled": 1479463552,
    "dateVisited": 0,
    "notes": "do blood test"
  },
  {
    "appointmentUUID": "987b09ee-543b-4dbe-99da-fd9c99202eee",
    "patientUUID": "6e894f6b-cbf6-4703-ad4f-bd93126450cb",
    "doctorUUID": "57c7aea1-9fea-422d-ae35-dbf8ce5f5dda",
    "patientName": "Joey Kapow",
    "dateScheduled": 1479463552,
    "dateVisited": 0,
    "notes": "check blood pressure"
  },
  {
    "appointmentUUID": "8bbd9f18-829b-4011-a451-df571b369796",
    "patientUUID": "e572fe98-4662-47f7-930c-cf4f7d13e26e",
    "doctorUUID": "57c7aea1-9fea-422d-ae35-dbf8ce5f5dda",
    "patientName": "Lamar Odom",
    "dateScheduled": 0,
    "dateVisited": 1479463552,
    "notes": "had difficulty breathing"
  }
]
```

-------------------------------------------------------

GET /doctors/doctoruuid/{doctoruuid}

**Retrieves a doctor profile**

Response:

HTTP 200 Found

```json
{
  "doctorUUID": "556d9f18-829b-4011-a451-df571b369111",
  "name": "Anoosh Gilliam",
  "phoneNumber": "555-222-1111",
  "primaryFacility": "address",
  "primarySpeciality": "Liver",
  "gender": "Female"
}
```
-------------------------------------------------------

GET /doctors

**Retrieves a list of all a clinics doctors**

Response:

HTTP 200 Found

```json
[
{
  "doctorUUID": "556d9f18-829b-4011-a451-df571b369111",
  "name": "Anoosh Gilliam",
  "phoneNumber": "555-222-1111",
  "primaryFacility": "address",
  "primarySpeciality": "Liver",
  "gender": "Female"
}
]
```
-------------------------------------------------------


POST /doctors

**Create a new doctor profile**

Request:

```json
{
  "name": "Doctor Name",
  "phoneNumber": "111-333-2222",
  "primaryFacility": "address",
  "primarySpeciality": "Specialty",
  "gender": "Male"
}
```

Response:

HTTP 201 Created

```json
{
  "code": 201,
  "message": "Doctor entry successfully created."
}
```
-------------------------------------------------------

POST /login

**Validates user credentials and returns userUUID**
**Requires using form body input (postman) or x-www-formurlencoded**
Request:

```
Form Data:
username:
password:
```

Responses:

HTTP 401 Unauthorized

```json
{
  "code": "401",
  "message": "Incorrect username or password"
}
```

HTTP 200 Found

```json
{
  "name": "Wolverine",
  "role": "Doctor",
  "userUUID": "556d9f18-829b-4011-a451-df571b369111"
}
```
-------------------------------------------------------

GET /users/useruuid/{useruuid}

**Get users basic information**

Responses:

HTTP 200 Found

```json
{
  "name": "Wolverine",
  "role": "Doctor",
  "userUUID": "556d9f18-829b-4011-a451-df571b369111"
}
```
-------------------------------------------------------


POST /users

**Create a new user entry**

Request:

```json
  "userName": "wolverine@xmen.ca",
  "passWord": "xmen",
  "role": "Doctor",
  "name": "Wolverine",
  "verification": "verificationKey"
}
```

Response:

HTTP 201 Created

```json
{
  "userUUID": "556d9f18-829b-4011-a451-df571b369111"
}
```
-------------------------------------------------------

PUT /patients

**Update an new user entry**

Request:

```json
{
  "patientUUID": "556d9f18-829b-4011-a451-df571b369111",
  "address": "5698 Cedar Avenue, San Francisco, California",
  "bloodType": "B-Positive",
  "dateOfBirth": 191289600,
  "emergencyContact": "415-555-8271",
  "gender": "F",
  "medicalNumber": "1234567890",
  "name": "Kelly Lai",
  "notes": "Accompanied by guide dog, ensure patient's area is wheelchair-friendly",
  "phoneNumber": "483-555-5123"
}
```

Response:

HTTP 200 OK

```json
{
  "code": 200,
  "message": "Patient entry successfully updated."
}
```

HTTP 400 OK

```json
{
  "code": 400,
  "message": "Error Occured: Patient not updated."
}
```
-------------------------------------------------------

POST /prescription

**Create a new prescription for a patient**

Request:

```json
{
  "patientUUID": 556d9f18-829b-4011-a451-df571b369111,
  "doctorUUID": 40119f18-829b-4011-a451-b369111df571.
  "doctor,omitempty": "Dr Ramoray",
  "drug": "Drug Name",
  "startDate": 191289600,
  "endDate": 191389600,
  "instructions,omitempty": "take twice daily"
}
```

Response:

HTTP 201 Created

```json
{
  "code": 201,
  "message": "Prescription entry successfully created."
}
```
-------------------------------------------------------
GET /prescriptions/patientuuid/{patientuuid}

**Retrieves a list of a patients prescriptions **

Response:

HTTP 200 Found

```json
[
{
  "patientUUID": 556d9f18-829b-4011-a451-df571b369111,
  "doctorUUID": 40119f18-829b-4011-a451-b369111df571.
  "prescriptionUUID": 57c7aea1-9fea-422d-ae35-dbf8ce5f5dda
  "doctor,omitempty": "Dr Ramoray",
  "drug": "Drug Name",
  "startDate": 191289600,
  "endDate": 191389600,
  "instructions,omitempty": "take twice daily"
}
]
```
-------------------------------------------------------


POST /notifications

**Create a new notificatoins for a doctor**

Request:

```json
{
"date":1488254862,
"message":"Hey man?",
"senderUUID":"20a5e81c-399f-4777-8bea-9c1fc2388f37",
"receiverUUID":"4498720b-0491-424f-8e52-6e13bd33da71",
"senderName":"Cyclops"
}
```

Response:

HTTP 201 Created

```json
{
  "code": 201,
  "message": "Notification entry successfully created."
}
```
-------------------------------------------------------
GET /notifications/doctoruuid/{doctoruuid}

**Retrieves a list of a doctor's notifications**

Response:

HTTP 200 Found

```json
[
  {
    "date": 1488254862,
    "message": "Have you seen Jean?",
    "receiverUUID": "4498720b-0491-424f-8e52-6e13bd33da71",
    "senderName": "Cyclops",
    "senderUUID": "20a5e81c-399f-4777-8bea-9c1fc2388f37"
  },
  {
    "date": 1388254862,
    "message": "Hey man?",
    "receiverUUID": "4498720b-0491-424f-8e52-6e13bd33da71",
    "senderName": "Cyclops",
    "senderUUID": "20a5e81c-399f-4777-8bea-9c1fc2388f37"
  }
]
```
-------------------------------------------------------
POST /documents

**Upload a new document associated with a patient**
**Send content as base64 encoded string of upload file**

Request:
```Form
------WebKitFormBoundaryAXbAxCjAnAVZ9VYz
Content-Disposition: form-data; name="dateUploaded"

1491035314
------WebKitFormBoundaryAXbAxCjAnAVZ9VYz
Content-Disposition: form-data; name="filename"

pdf.pdf
------WebKitFormBoundaryAXbAxCjAnAVZ9VYz
Content-Disposition: form-data; name="patientUUID"

36b95ee0-3742-42a1-a521-ecbb2528e2a4
------WebKitFormBoundaryAXbAxCjAnAVZ9VYz
Content-Disposition: form-data; name="file"; filename="pdf.pdf"
Content-Type: application/pdf


------WebKitFormBoundaryAXbAxCjAnAVZ9VYz--
```

Response:

HTTP 201 Created

```json
{
  "code": 201,
  "message": "Document entry successfully created."
}
```
-------------------------------------------------------
GET /documents/documentuuid/{documentuuid}

**Downloads a document**

Response:

HTTP 200 Found (Downloads the requested file)

-------------------------------------------------------

GET /documents/patientuuid/{patientuuid}

**Retrieves a list of a patients's documents**

Response:

HTTP 200 Found

```json
[
  {
    "documentUUID": "4498720b-0491-424f-8e52-6e13bd33da71",
    "patientUUID": "20a5e81c-399f-4777-8bea-9c1fc2388f37",
    "filename": "fileName.pdf"
    "dateUploaded": 1488254862
  }
]
```
-------------------------------------------------------
