package main

import (
"fmt"
	"github.com/go-pg/migrations"
	_ "github.com/lib/pq"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		fmt.Println("running initial migration")
		_, err := db.Exec(`
CREATE TYPE job_status AS ENUM ('filling', 'underway', 'completed', 'cancelled');
CREATE TYPE contractor_job_status AS ENUM ('invited', 'requesting', 'declined', 'approved');
CREATE TYPE user_roles AS ENUM ('manager', 'contractor', 'admin');

CREATE TABLE users(
    id SERIAL UNIQUE PRIMARY KEY, 
    email varchar(100) UNIQUE NOT NULL,
    password_hash varchar(100) NOT NULL,
	role user_roles NOT NULL
);

CREATE INDEX IndexUsersEmail
ON users (email);

CREATE TABLE auth_tokens(
    id SERIAL UNIQUE PRIMARY KEY, 
    token_id varchar(100) UNIQUE NOT NULL,
    user_id INT NOT NULL
);

CREATE INDEX IndexAuthTokens
ON auth_tokens (user_id);

CREATE TABLE contractors(
    id SERIAL UNIQUE PRIMARY KEY, 
    name   varchar(50) NOT NULL,
    charge_rate varchar(7) NOT NULL,
    email      varchar(50) NOT NULL,
    enabled    BOOLEAN NOT NULL DEFAULT TRUE,
    notes varchar(1000) ,
    phone      varchar(15)      NOT NULL,
	company_id INT NOT NULL,
	available bool NOT NULL DEFAULT TRUE,
	due_back TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IndexContractorsChargeRate
ON contractors (charge_rate);
CREATE INDEX IndexContractorsEnabled
ON contractors (company_id, enabled);

CREATE TABLE managers(
    id SERIAL UNIQUE PRIMARY KEY, 
    name varchar(50) NOT NULL,
    email varchar(50) NOT NULL,
    phone varchar(15) NOT NULL,
    company_id int NOT NULL
);
CREATE INDEX IndexCompanyId
ON managers (company_id);

CREATE TABLE jobs(
    id SERIAL UNIQUE PRIMARY KEY, 
    name   varchar(50) NOT NULL,
	effort varchar(50) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    status job_status NOT NULL,
    description varchar(300) NOT NULL,
    manager_id Int NOT NULL
);
CREATE INDEX IndexJobManagerId
ON jobs (manager_id);
CREATE INDEX IndexJobStatusManagerId
ON jobs (status, manager_id);

CREATE TABLE companies(
    id SERIAL UNIQUE PRIMARY KEY, 
    name varchar(50) NOT NULL
);

CREATE TABLE skills(
    id SERIAL UNIQUE PRIMARY KEY, 
    name   varCHAR(50) NOT NULL
);

CREATE TABLE contractor_skills(
	id SERIAL UNIQUE PRIMARY KEY,
    contractor_id INT NOT NULL, 
    skill_id INT NOT NULL
);
CREATE INDEX IndexContractorSkills
ON contractor_skills (skill_id, contractor_id);
CREATE INDEX IndexContractorSkillsContractor
ON contractor_skills (contractor_id);

CREATE TABLE contractor_jobs(
    id SERIAL UNIQUE PRIMARY KEY, 
    contractor_id INT Not NUll,
    status contractor_job_status NOT NULL,
	state_seen bool NOT NULL DEFAULT FALSE,
	job_id INT NOT NULL
);
CREATE INDEX IndexContractorJobsContractorId
ON contractor_jobs (contractor_id, status);
CREATE INDEX IndexContractorJobsContractorIdStateSeen
ON contractor_jobs (contractor_id, status, state_seen);

CREATE TABLE company_skills(
    id SERIAL UNIQUE PRIMARY KEY, 
    company_id INT NOT NULL,
    skill_id INT NOT NULL
);
CREATE INDEX IndexCompanySkills
ON company_skills (company_id, skill_id);
`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("removing everything!!!")
		_, err := db.Exec(`
DROP TABLE contractors;
DROP TABLE managers;
DROP TABLE jobs;
DROP TABLE companies;
DROP TABLE skills;
DROP TABLE contractor_skills;
DROP TABLE contractor_jobs;
DROP TABLE company_skills;
DROP TYPE job_status;
DROP TYPE contractor_job_status;
		`)
		return err
	})
}
