package models

import (
	"database/sql"
	"strings"
)

type ContractorJob struct {
	ID           int    `json:"id" binding:"required"`
	ContractorID int    `json:"contractor_id" binding:"required"`
	Status       string `json:"status" binding:"required"`
	StateSeen    bool   `json:"state_seen" binding:"required"`
	JobID        int    `json:"job_id" binding:"required"`
}

func (c *ContractorJob) GetContractorJob(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM contractor_jobs WHERE contractor_id=$1 AND job_id=$2",
		c.ContractorID, c.JobID).Scan(&c.ID, &c.ContractorID, &c.Status, &c.StateSeen, &c.JobID)
}

func (c *ContractorJob) UpdateContractorJob(db *sql.DB) error {
	_, err := db.Exec("UPDATE contractor_jobs SET status=$1, state_seen=$2 WHERE contractor_id=$3 AND job_id=$4",
			c.Status, c.StateSeen, c.ContractorID, c.JobID)

	return err
}

func (c *ContractorJob) DeleteContractorJob(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM contractor_jobs WHERE contractor_id=$1 AND job_id=$2", c.ContractorID, c.JobID)

	return err
}

func (c *ContractorJob) CreateContractorJob(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO contractor_jobs(contractor_id, status, state_seen, job_id) "+
		"VALUES($1, $2, $3, $4) RETURNING id", c.ContractorID, c.Status, c.StateSeen, c.JobID).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetContractorJobs(db *sql.DB, contractorId string, statuses []string) ([]ContractorJob, error) {
	query := QueryBuilder{}
	var params []interface{}
	params = append(params, contractorId)
	query = query.AddQueryString("SELECT * FROM contractor_jobs", false).AddWhereClause("contractor_id=$1", params)

	if len(statuses[0]) != 0 {
		params = nil
		params = append(params, strings.Join(statuses, "','"))
		query = query.AddWhereClause("status IN ('$1')", params)
	}

	rows, err := query.Get(db)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	contractorJobs := make([]ContractorJob, 0)
	for rows.Next() {
		var c ContractorJob
		if err := rows.Scan(&c.ID, &c.ContractorID, &c.Status, &c.StateSeen, &c.JobID); err != nil {
			return nil, err
		}
		contractorJobs = append(contractorJobs, c)
	}

	return contractorJobs, nil
}
