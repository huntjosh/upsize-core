package models

import (
	"database/sql"
)

type Skill struct {
	ID   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

func (s *Skill) GetSkill(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM skills WHERE id=$1", s.ID).Scan(&s.ID, &s.Name)
}

func (s *Skill) UpdateSkill(db *sql.DB) error {
	_, err := db.Exec("UPDATE skills SET name=$1 WHERE id=$2", s.Name, s.ID)

	return err
}

func (s *Skill) DeleteSkill(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM skills WHERE id=$1", s.ID)

	return err
}

func (s *Skill) CreateSkill(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO skills(name) VALUES($1) RETURNING id", s.Name).Scan(&s.ID)

	if err != nil {
		return err
	}

	return nil
}

func GetSkills(db *sql.DB, start, count int) ([]Skill, error) {
	rows, err := db.Query("SELECT * FROM skills LIMIT $1 OFFSET $2", count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	skills := make([]Skill, 0)
	for rows.Next() {
		var s Skill
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		skills = append(skills, s)
	}

	return skills, nil
}
