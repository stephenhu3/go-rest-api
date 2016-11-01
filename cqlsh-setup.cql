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


CREATE TABLE completed_appointments (
	appointmentUuid uuid,
	patientUuid uuid,
	doctorUuid uuid,
	dateVisited timestamp,
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


CREATE TABLE scheduled_appointments (
	appointmentUuid uuid,
	patientUuid uuid,
	doctorUuid uuid,
	dateScheduled timestamp,
	PRIMARY KEY (appointmentUuid)
);
CREATE INDEX patientUuid ON emr.scheduled_appointments (patientUuid);
CREATE INDEX doctorUuid ON emr.scheduled_appointments (dateVisited);
CREATE INDEX dateScheduled ON emr.scheduled_appointments (dateScheduled);