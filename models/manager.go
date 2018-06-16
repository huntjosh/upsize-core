package models

import (
	"database/sql"
	"strconv"
)

type Manager struct {
	ID    int     `json:"id" binding:"required"`
	Name  string  `json:"name" binding:"required"`
	Email  string  `json:"email" binding:"required"`
	Phone  string  `json:"phone" binding:"required"`
	CompanyID  int  `json:"company_id" binding:"required"`
}

func (c *Company) GetManagerCompany(db *sql.DB, managerId string) error {
	err := db.QueryRow("SELECT companies.* FROM managers JOIN companies ON managers.company_id = companies.id WHERE managers.id=$1",
		managerId).Scan(&c.ID, &c.Name)

	return err
}

func (m *Manager) GetManager(db *sql.DB) error {
	err := db.QueryRow("SELECT * FROM managers WHERE id=$1", m.ID).Scan(&m.ID, &m.Name, &m.Email, &m.Phone, &m.CompanyID)

	return err
}

func (m *Manager) UpdateManager(db *sql.DB) error {
	params := make([]interface{}, 0)
	params = append(params, m.ID)
	paramCount := 1
	valuesString := ""

	if m.Name != "" {
		paramCount++
		params = append(params, m.Name)
		valuesString += " name=$" + strconv.Itoa(paramCount)
	}

	if m.Phone != "" {
		paramCount++
		params = append(params, m.Phone)
		if valuesString != "" {
			valuesString += ","
		}
		valuesString += " phone=$" + strconv.Itoa(paramCount)
	}

	if m.Email != "" {
		paramCount++
		params = append(params, m.Email)
		if valuesString != "" {
			valuesString += ","
		}
		valuesString += " email=$" + strconv.Itoa(paramCount)
	}

	_, err :=
		db.Exec("UPDATE managers SET"+valuesString+" WHERE id=$1",
			params...)

	return err
}

func (m *Manager) DeleteManager(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM managers WHERE id=$1", m.ID)

	return err
}

func (m *Manager) CreateManager(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO managers(name, email, phone, company_id) VALUES($1, $2, $3, $4) RETURNING id",
		m.Name, m.Email, m.Phone, m.CompanyID).Scan(&m.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetManagers(db *sql.DB, companyId string) ([]Manager, error) {
	var rows *sql.Rows
	var err error

	if companyId != "" {
		rows, err = db.Query(
			"SELECT * FROM managers WHERE company_id=$1",
			companyId)
	} else {
		rows, err = db.Query(
			"SELECT * FROM managers")
	}


	if err != nil {
		return nil, err
	}

	defer rows.Close()

	managers := make([]Manager, 0)
	for rows.Next() {
		var m Manager
		if err := rows.Scan(&m.ID, &m.Name, &m.Email, &m.Phone, &m.CompanyID); err != nil {
			return nil, err
		}
		managers = append(managers, m)
	}

	return managers, nil
}