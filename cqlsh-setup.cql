DROP KEYSPACE IF EXISTS emr;

CREATE KEYSPACE IF NOT EXISTS emr
	WITH replication = {
	  'class': 'SimpleStrategy',
	  'replication_factor': '1'
	};

USE emr;

CREATE TABLE patients (
	patientUUID uuid,
	dateOfBirth int,
	gender text,
	name text,
	medicalNumber text,
	bloodType text,
	emergencyContact text,
	phone text,
	address text,
	notes text,
	PRIMARY KEY (patientUUID)
);


CREATE TABLE completedAppointments (
	appointmentUUID uuid,
	patientUUID uuid,
	doctorUUID uuid,
	dateVisited int,
	breathingRate int,
	heartRate int,
	bloodOxygenLevel int,
	bloodPressure int,
	notes text,
	PRIMARY KEY (appointmentUUID)
);
CREATE INDEX completedAppointmentPatientUUID ON emr.completedAppointments (patientUUID);
CREATE INDEX completedAppointmentDoctorUUID ON emr.completedAppointments (doctorUUID);
CREATE INDEX completedAppointmentDateVisited ON emr.completedAppointments (dateVisited);


CREATE TABLE futureAppointments (
	appointmentUUID uuid,
	patientUUID uuid,
	doctorUUID uuid,
	dateScheduled int,
	notes text,
	PRIMARY KEY (appointmentUUID)
);
CREATE INDEX futureAppointmentPatientUUID ON emr.futureAppointments (patientUUID);
CREATE INDEX futureAppointmentDoctorUUID ON emr.futureAppointments (doctorUUID);
CREATE INDEX futureAppointmentDateScheduled ON emr.futureAppointments (dateScheduled);


CREATE TABLE doctors (
	doctorUUID uuid,
	name text,
	phone text,
	primaryFacility text,
	primarySpecialty text,
	gender text,
	PRIMARY KEY (doctorUUID)
);

CREATE TABLE users (
	username text,
	salt text,
	saltedHash text,
	userUUID uuid,
	role text,
	name text,
	PRIMARY KEY (username)
);
CREATE INDEX usersUserUUID ON emr.users (userUUID);

CREATE TABLE prescriptions (
	patientUUID uuid,
	prescriptionUUID uuid,
	doctorUUID uuid,
	doctorName text,
	drug text,
	startDate int,
	endDate int,
	instructions text,
	PRIMARY KEY (patientUUID, endDate, prescriptionUUID)
);

CREATE TABLE notifications (
	date int,
	message text,
	notificationuuid uuid,
	receiverUUID uuid,
	senderName text,
	senderUUID uuid,
	PRIMARY KEY (receiverUUID, date, notificationuuid)
) WITH CLUSTERING ORDER BY (date DESC);

CREATE TABLE documents (
	documentUUID uuid,
	patientUUID uuid,
	filename text,
	content blob,
	dateUploaded int,
	PRIMARY KEY (documentUUID)
);
CREATE INDEX documentsPatientUUID ON emr.documents (patientUUID);
