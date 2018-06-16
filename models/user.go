package models

import (
	"database/sql"
	"strconv"
	"fmt"
)

type User struct {
	ID           int    `json:"id" binding:"required"`
	Email        string `json:"email" binding:"required"`
	PasswordHash string `json:"password_hash" binding:"required"`
	Role         string `json:"role" binding:"required"`
}

func (u *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE email=$1", u.Email)

	return err
}

func (u *User) UpdateUser(db *sql.DB) error {
	params := make([]interface{}, 0)
	params = append(params, u.Email)
	paramCount := 1
	valuesString := ""

	if u.PasswordHash == "" && u.Role == "" {
		return nil
	}

	if u.PasswordHash != "" {
		paramCount++
		params = append(params, u.PasswordHash)
		valuesString += " password_hash=$" + strconv.Itoa(paramCount)
	}

	if u.Role != "" {
		paramCount++
		params = append(params, u.Role)
		if valuesString != "" {
			valuesString += ","
		}
		valuesString += " role=$" + strconv.Itoa(paramCount)
	}

	_, err :=
		db.Exec("UPDATE users SET"+valuesString+" WHERE email=$1",
			params...)

	return err
}

func (u *User) CreateUser(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO users(email, password_hash, role) VALUES($1, $2, $3) RETURNING id",
		u.Email, u.PasswordHash, u.Role).Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func (u *User) GetUser(db *sql.DB) error {
	err := db.QueryRow("SELECT * FROM users WHERE email=$1", u.Email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func (u *User) GetUserNoPassword(db *sql.DB) error {
	err := db.QueryRow("SELECT id, email, role FROM users WHERE email=$1", u.Email).Scan(&u.ID, &u.Email, &u.Role)
	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}

func GetCompanyIDFromEmail(db *sql.DB, email, role string) int {
	var companyId int
	err := db.QueryRow("SELECT "+role+"s.company_id FROM users JOIN "+role+"s ON users.email = "+role+"s.email WHERE users.email=$1", email).Scan(&companyId)
	if err != nil {
		return 0
	}
	return companyId
}

func GetCompanyIDFromID(db *sql.DB, id, role string) int {
	var companyId int
	err := db.QueryRow("SELECT "+role+"s.company_id FROM "+role+"s WHERE id=$1", id).Scan(&companyId)
	if err != nil {
		return 0
	}
	return companyId
}

func GetId(db *sql.DB, email, role string) int {
	var id int
	err := db.QueryRow("SELECT "+role+"s.id FROM users JOIN "+role+"s ON users.email = "+role+"s.email WHERE users.email=$1", email).Scan(&id)
	if err != nil {
		return 0
	}
	return id
}