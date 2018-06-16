package models

import (
	"database/sql"
)

type CompanySkill struct {
	ID        int `json:"id" binding:"required"`
	SkillID   int `json:"skill_id" binding:"required"`
	CompanyID int `json:"company_id" binding:"required"`
}

func (cs *CompanySkill) DeleteCompanySkill(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM company_skills WHERE company_id=$1 AND skill_id=$2", cs.CompanyID, cs.SkillID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CompanySkill) CreateCompanySkill(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO company_skills(skill_id, company_id) VALUES($1, $2) RETURNING id",
		cs.SkillID, cs.CompanyID).Scan(&cs.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Skill) GetCompanySkill(db *sql.DB, companyId string, skillId string) error {
	return db.QueryRow("SELECT skills.* FROM skills JOIN company_skills ON skills.id = company_skills.skill_id "+
		"WHERE company_skills.company_id=$1 AND company_skills.skill_id=$2", companyId, skillId).Scan(&s.ID, &s.Name)
}

func GetCompanySkills(db *sql.DB, companyId string) ([]Skill, error) {
	rows, err := db.Query("SELECT skills.* FROM skills JOIN company_skills ON skills.id = company_skills.skill_id "+
		"WHERE company_skills.company_id=$1", companyId)
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
