package models

import (
	"database/sql"
	"time"
	"github.com/lib/pq"
)

type Contractor struct {
	ID    int     `json:"id" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	ChargeRate  string  `json:"charge_rate" binding:"required"`
	Email  string  `json:"email" binding:"required"`
	Enabled  bool  `json:"enabled" binding:"required"`
	Notes  string  `json:"notes"`
	Phone  string  `json:"phone" binding:"required"`
	CompanyID  int  `json:"company_id" binding:"required"`
	Available  bool  `json:"available" binding:"required"`
	DueBack  time.Time  `json:"due_back"`
}

func (c *Company) GetContractorCompany(db *sql.DB, contractorId string) error {
	err := db.QueryRow("SELECT companies.* FROM contractors JOIN companies ON contractors.company_id = companies.id WHERE contractors.id=$1",
		contractorId).Scan(&c.ID, &c.Name)
	return err
}

func (c *Contractor) GetContractor(db *sql.DB) error {
	var dueBack pq.NullTime
	var notes sql.NullString

	err := db.QueryRow("SELECT * FROM contractors WHERE id=$1",
		c.ID).Scan(&c.ID, &c.Name, &c.ChargeRate, &c.Email, &c.Enabled, &notes, &c.Phone, &c.CompanyID, &c.Available, &dueBack)
	if notes.Valid {
		c.Notes = notes.String
	}
	if dueBack.Valid {
		c.DueBack = dueBack.Time
	}

	return err
}

func (c *Contractor) UpdateContractor(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE contractors SET name=$1, charge_rate=$2, email=$3, enabled=$4, notes=$5, phone=$6, company_id=$7, available=$8, due_back=$9 WHERE id=$10",
			c.Name, c.ChargeRate, c.Email, c.Enabled, c.Notes, c.Phone, c.CompanyID, c.Available, c.DueBack, c.ID)

	return err
}

func (c *Contractor) DeleteContractor(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM contractors WHERE id=$1", c.ID)

	return err
}

func (c *Contractor) CreateContractor(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO contractors(name, charge_rate, email, enabled, notes, phone, company_id, available, due_back) "+
			"VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", c.Name, c.ChargeRate, c.Email, c.Enabled, c.Notes,
		c.Phone, c.CompanyID, c.Available, c.DueBack).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetContractors(db *sql.DB, companyId string) ([]Contractor, error) {
	query := QueryBuilder{}
	query = query.AddQueryString("SELECT * FROM contractors", false)

	if companyId != "" {
		query = AddCompanyWhereClause(query, companyId)
	}

	rows, err := query.Get(db)
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

func MapRowToContractor(rows *sql.Rows) (Contractor, error) {
	var c Contractor
	var dueBack pq.NullTime
	var notes sql.NullString

	if err := rows.Scan(&c.ID, &c.Name, &c.ChargeRate, &c.Email, &c.Enabled, &notes, &c.Phone, &c.CompanyID,
		&c.Available, &dueBack); err != nil {
		return Contractor{}, err
	}

	if dueBack.Valid {
		c.DueBack = dueBack.Time
	}

	if notes.Valid {
		c.Notes = notes.String
	}

	return c, nil
}

func GetCompanyContractors(db *sql.DB, companyId string) ([]Contractor, error) {
	rows, err := db.Query(
		"SELECT contractors.* FROM contractors JOIN companies ON contractors.company_id = companies.id WHERE companies.id=$1 ",
		companyId)

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