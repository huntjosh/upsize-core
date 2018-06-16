package models

import (
	"database/sql"
)

type Company struct {
	ID    int     `json:"id" binding:"required"`
	Name  string  `json:"name" binding:"required"`
}

func (c *Company) GetCompany(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM companies WHERE id=$1",
		c.ID).Scan(&c.ID, &c.Name)
}

func (c *Company) UpdateCompany(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE companies SET name=$1 WHERE id=$2",
			c.Name, c.ID)

	return err
}

func (c *Company) DeleteCompany(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM companies WHERE id=$1", c.ID)

	return err
}

func (c *Company) CreateCompany(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO companies(name) VALUES($1) RETURNING id",
		c.Name).Scan(&c.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetCompanies(db *sql.DB) ([]Company, error) {
	rows, err := db.Query("SELECT * FROM companies")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	companies := make([]Company, 0)
	for rows.Next() {
		var c Company
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}

	return companies, nil
}

func AddCompanyWhereClause(q QueryBuilder, companyId string) QueryBuilder {
	var params []interface{}
	params = append(params, companyId)
	return q.AddWhereClause("company_id=$1", params)
}