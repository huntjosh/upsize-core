package models

import (
	"database/sql"
	"time"
	"github.com/lib/pq"
	"strconv"
)

type Job struct {
	ID          int       `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Effort      string    `json:"effort" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date"`
	Status      string    `json:"status" binding:"required"`
	Description string    `json:"description" binding:"required"`
	ManagerID   int       `json:"manager_id" binding:"required"`
}

func (j *Job) GetJob(db *sql.DB) error {
	var endDate pq.NullTime
	err := db.QueryRow("SELECT id, name, effort, start_date, end_date, status, description, manager_id FROM jobs WHERE id=$1",
		j.ID).Scan(&j.ID, &j.Name, &j.Effort, &j.StartDate, &endDate, &j.Status, &j.Description, &j.ManagerID)

	if endDate.Valid {
		j.EndDate = endDate.Time
	}

	return err
}

func (j *Job) UpdateJob(db *sql.DB) error {
	_, err := db.Exec("UPDATE jobs SET name=$1, effort=$2, start_date=$3, end_date=$4, status=$5, "+
		"description=$6, manager_id=$7 WHERE id=$8", j.Name, j.Effort, j.StartDate, j.EndDate, j.Status, j.Description,
		j.ManagerID, j.ID)

	return err
}

func (j *Job) DeleteJob(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM jobs WHERE id=$1", j.ID)

	return err
}

func (j *Job) CreateJob(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO jobs(name, effort, start_date, end_date, status, description, manager_id) "+
		"VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id", j.Name, j.Effort, j.StartDate, j.EndDate, j.Status,
		j.Description, j.ManagerID).Scan(&j.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetJobs(db *sql.DB, companyId string) ([]Job, error) {
	query := QueryBuilder{}
	query = query.AddQueryString("SELECT jobs.* FROM jobs JOIN managers on jobs.manager_id = managers.id ", false)
	if companyId != "" {
		var params []interface{}
		params = append(params, companyId)
		query = query.AddWhereClause("managers.company_id=$1", params)
	}

	rows, err := query.Get(db)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jobs := make([]Job, 0)
	for rows.Next() {
		j, err := MapRowToJob(rows)
		if err == nil {
			jobs = append(jobs, j)
		}
	}

	return jobs, nil
}

func GetCompanyJobs(db *sql.DB, companyId string, statuses []string) ([]Job, error) {
	whereClause, params := idStatusParams(companyId, statuses)
	rows, err := db.Query("SELECT jobs.* FROM jobs JOIN managers ON jobs.manager_id = managers.id JOIN companies ON "+
			"companies.id = managers.company_id WHERE companies.id=$1 "+ whereClause, params...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jobs := make([]Job, 0)
	for rows.Next() {
		j, err := MapRowToJob(rows)
		if err == nil {
			jobs = append(jobs, j)
		}
	}

	return jobs, nil
}

func GetManagerJobs(db *sql.DB, managerId string, statuses []string) ([]Job, error) {
	whereClause, params := idStatusParams(managerId, statuses)
	rows, err := db.Query("SELECT jobs.* FROM jobs WHERE manager_id = $1 "+whereClause, params...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jobs := make([]Job, 0)
	for rows.Next() {
		j, err := MapRowToJob(rows)
		if err == nil {
			jobs = append(jobs, j)
		}
	}

	return jobs, nil
}

func GetJobContractors(db *sql.DB, jobId string) ([]Contractor, error) {
	rows, err := db.Query("SELECT contractors.* FROM contractor_jobs JOIN contractors "+
		"ON contractor_jobs.contractor_id = contractors.id WHERE contractor_jobs.job_id=$1", jobId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	contractors := make([]Contractor, 0)
	for rows.Next() {
		c, err := MapRowToContractor(rows)
		if err == nil {
			contractors = append(contractors, c)
		}
	}

	return contractors, nil
}

func MapRowToJob(rows *sql.Rows) (Job, error) {
	var j Job
	var endDate pq.NullTime

	if err := rows.Scan(&j.ID, &j.Name, &j.Effort, &j.StartDate, &endDate, &j.Status, &j.Description, &j.ManagerID); err != nil {
		return Job{}, err
	}

	if endDate.Valid {
		j.EndDate = endDate.Time
	}

	return j, nil
}

func idStatusParams(id string, statuses []string) (string, []interface{}) {
	params := make([]interface{}, 0)
	params = append(params, id)
	whereClause := ""
	if len(statuses[0]) != 0 { // Because the slice is initialized with 1 value
		statusParams := ""
		for i := 0; i < len(statuses); i++ {
			if i > 0 {
				statusParams += ", "
			}
			statusParams += "$" + strconv.Itoa(2+i) // First param is for the manager id
			params = append(params, statuses[i])
		}
		whereClause = " AND status IN (" + statusParams + ")"
	}

	return whereClause, params
}

func GetCompanyFromJobID(db *sql.DB, jobId string) int {
	var companyId int
	err := db.QueryRow("SELECT managers.company_id FROM jobs JOIN managers on jobs.manager_id = managers.id "+
		"WHERE jobs.id=$1", jobId).Scan(&companyId)
	if err != nil {
		return 0
	}

	return companyId
}
