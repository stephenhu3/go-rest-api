CREATE KEYSPACE IF NOT EXISTS emr
	WITH replication = {
	  'class': 'SimpleStrategy',
	  'replication_factor': '1'
	};

USE emr;

CREATE TABLE patients (
	patientUUID uuid,
	age int,
	dateOfBirth date,
	dateOfDeath date,
	gender text,
	name text,
	medicalNumber text,
	ethnicity text,
	emergencyContact text,
	phoneNumber text,
	address text,
	notes: text,
	PRIMARY KEY (patientUUID, ethnicity)
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
CREATE INDEX completedAppointmentdoctorUUID ON emr.completedAppointments (doctorUUID);
CREATE INDEX completedAppointmentdateVisited ON emr.completedAppointments (dateVisited);


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
CREATE INDEX dateScheduled ON emr.futureAppointments (dateScheduled);
