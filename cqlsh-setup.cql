CREATE KEYSPACE emr WITH replication = {
  'class': 'SimpleStrategy',
  'replication_factor': '1'
};

USE emr;

CREATE TABLE patients (
	patientUUID uuid,
	age int,
	gender text,
	name text,
	insuranceNumber text,
	PRIMARY KEY (patientUuid)
);


CREATE TABLE completedAppointments (
	appointmentUuid uuid,
	patientUuid uuid,
	doctorUuid uuid,
	dateVisited int,
	breathingRate int,
	heartRate int,
	bloodOxygenLevel int,
	bloodPressure int,
	notes text,
	PRIMARY KEY (appointmentUuid)
);
CREATE INDEX patientUuid ON emr.completedAppointments (patientUuid);
CREATE INDEX doctorUuid ON emr.completedAppointments (doctorUuid);
CREATE INDEX dateVisited ON emr.completedAppointments (dateVisited);


CREATE TABLE futureAppointments (
	appointmentUuid uuid,
	patientUuid uuid,
	doctorUuid uuid,
	dateScheduled int,
	notes text,
	PRIMARY KEY (appointmentUuid)
);
CREATE INDEX schedPatientUuid ON emr.futureAppointments (patientUuid);
CREATE INDEX schedDoctorUuid ON emr.futureAppointments (doctorUuid);
CREATE INDEX dateScheduled ON emr.futureAppointments (dateScheduled);
