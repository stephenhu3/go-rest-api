CREATE KEYSPACE emr WITH replication = {
  'class': 'SimpleStrategy',
  'replication_factor': '1'
};

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
CREATE INDEX patientUuid ON emr.completed_appointments (patientUuid);
CREATE INDEX doctorUuid ON emr.completed_appointments (doctorUuid);
CREATE INDEX dateVisited ON emr.completed_appointments (dateVisited);


CREATE TABLE futureAppointments (
	appointmentUuid uuid,
	patientUuid uuid,
	doctorUuid uuid,
	dateScheduled int,
	notes text,
	PRIMARY KEY (appointmentUuid)
);
CREATE INDEX patientUuid ON emr.futureAppointments (patientUuid);
CREATE INDEX doctorUuid ON emr.futureAppointments (doctorUuid);
CREATE INDEX dateScheduled ON emr.futureAppointments (dateScheduled);