package tests

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"
	"upsizeAPI/models"
)

func TestCreateSkill(t *testing.T) {
	FreshDatabase()

	payload := []byte(`{"name":"best skill"}`)

	req, _ := http.NewRequest("PUT", "/skill", bytes.NewBuffer(payload))
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "best skill" {
		t.Errorf("Expected skill name to be 'best skill'. Got '%v'", m["name"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected skill ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetNonExistentSkill(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/skill/11", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Skill not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Skill not found'. Got '%s'", m["error"])
	}
}

func TestGetSkill(t *testing.T) {
	FreshDatabase()
	addSkills(1)

	req, _ := http.NewRequest("GET", "/skill/1", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addSkills(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO skills(name) VALUES($1)", "Skill "+strconv.Itoa(i))
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestUpdateSkill(t *testing.T) {
	FreshDatabase()
	addSkills(1)

	req, _ := http.NewRequest("GET", "/skill/1", nil)
	response := executeRequest(req, "manager")
	var originalSkill map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalSkill)

	payload := []byte(`{"name":"Skill - updated name"}`)

	req, _ = http.NewRequest("POST", "/skill/1", bytes.NewBuffer(payload))
	response = executeRequest(req, "admin")

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalSkill["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalSkill["id"], m["id"])
	}

	if m["name"] == originalSkill["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalSkill["name"], m["name"], m["name"])
	}
}

func TestDeleteSkill(t *testing.T) {
	FreshDatabase()
	addSkills(1)

	req, _ := http.NewRequest("GET", "/skill/1", nil)
	response := executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/skill/1", nil)
	response = executeRequest(req, "admin")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/skill/1", nil)
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetSkills(t *testing.T) {
	FreshDatabase()
	addSkills(4)

	req, _ := http.NewRequest("GET", "/skills", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var companies []models.Skill
	json.Unmarshal(response.Body.Bytes(), &companies)
	if len(companies) != 4 {
		t.Errorf("Expected skills retrieved to be 4, found " + strconv.Itoa(len(companies)))
	}
}
