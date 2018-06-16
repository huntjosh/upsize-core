package tests

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"
	"upsizeAPI/models"
)

func TestCreateCompanySkill(t *testing.T) {
	FreshDatabase()

	payload := []byte(`{"skill_id":5}`)

	req, _ := http.NewRequest("PUT", "/company/1/skill", bytes.NewBuffer(payload))
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["skill_id"] != 5.0 {
		t.Errorf("Expected skill_id to be '5'. Got '%v'", m["skill_id"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected company_skill ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetNonExistentCompanySkill(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/company/1/skill/11", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "company skill not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'company skill not found'. Got '%s'", m["error"])
	}
}

func TestGetCompanySkill(t *testing.T) {
	FreshDatabase()
	addCompanySkills(1)

	req, _ := http.NewRequest("GET", "/company/1/skill/1", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addCompanySkills(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO companies(name) VALUES($1)", "Company "+strconv.Itoa(i))
		if err != nil {
			panic(err.Error())
		}

		_, err = a.DB.Exec("INSERT INTO skills(name) VALUES($1)", "Skill "+strconv.Itoa(i))
		if err != nil {
			panic(err.Error())
		}

		_, err = a.DB.Exec("INSERT INTO company_skills(company_id,skill_id) VALUES($1, $2)", 1, i+1)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestDeleteCompanySkill(t *testing.T) {
	FreshDatabase()
	addCompanySkills(1)

	req, _ := http.NewRequest("GET", "/company/1/skill/1", nil)
	response := executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/company/1/skill/1", nil)
	response = executeRequest(req, "admin")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/company/1/skill/1", nil)
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetCompanySkills(t *testing.T) {
	FreshDatabase()
	addCompanySkills(5)

	req, _ := http.NewRequest("GET", "/company/1/skills", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var companySkills []models.Skill
	json.Unmarshal(response.Body.Bytes(), &companySkills)
	if len(companySkills) != 5 {
		t.Errorf("Expected company_skills retrieved to be 5, found " + strconv.Itoa(len(companySkills)))
	}
}
