package tests

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"
	"upsizeAPI/models"
)

func TestCreateContractor(t *testing.T) {
	FreshDatabase()

	payload := []byte(`
		{"name":"echo","charge_rate":"25","email":"echovvv@gmail.com","enabled":true,"notes":"123",
			"phone":"02040490234","company_id":1,"available":true,"due_back":"2018-01-08T04:05:06-01:00"}`)

	req, _ := http.NewRequest("PUT", "/contractor", bytes.NewBuffer(payload))
	response := executeRequest(req, "admin")

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "echo" {
		t.Errorf("Expected contractor name to be 'echo'. Got '%v'", m["name"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 2.0 {
		t.Errorf("Expected contractor ID to be '2'. Got '%v'", m["id"])
	}
}

func TestGetNonExistentContractor(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/contractor/11", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestGetContractor(t *testing.T) {
	FreshDatabase()
	addContractors(1)

	req, _ := http.NewRequest("GET", "/contractor/1", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addContractors(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO contractors(name, charge_rate, email, enabled, phone, company_id, available) VALUES($1, $2, $3, $4, $5, $6, $7)",
			"Contractor "+strconv.Itoa(i), "20", "j@gmail.com", true, "02040490234", 1, true)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestUpdateContractor(t *testing.T) {
	FreshDatabase()
	addContractors(1)

	req, _ := http.NewRequest("GET", "/contractor/1", nil)
	response := executeRequest(req, "manager")
	var originalContractor map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalContractor)

	payload := []byte(`{"name":"echo","charge_rate":"25","email":"echovvv@gmail.com","enabled":false,"notes":"123",
			"phone":"02049490234","company_id":2,"available":false,"due_back":"2018-01-08T04:05:06-01:00"}`)

	req, _ = http.NewRequest("POST", "/contractor/1", bytes.NewBuffer(payload))
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	attributes := []string{"name", "charge_rate", "email", "enabled", "notes", "phone", "company_id", "available", "due_back"}
	if m["id" ] != originalContractor["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalContractor["id"], m["id"])
	}

	for _, attributeName := range attributes {
		if m[attributeName] == originalContractor[attributeName] {
			t.Errorf("Expected the %s to change from '%v' to '%v'. Got '%v'", attributeName,
				originalContractor[attributeName], m[attributeName], m[attributeName])
		}
	}
}

func TestDeleteContractor(t *testing.T) {
	FreshDatabase()
	addContractors(1)

	req, _ := http.NewRequest("GET", "/contractor/1", nil)
	response := executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/contractor/1", nil)
	response = executeRequest(req, "admin")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/contractor/1", nil)
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestGetContractors(t *testing.T) {
	FreshDatabase()
	addContractors(2)

	req, _ := http.NewRequest("GET", "/contractors", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var contractors []models.Contractor
	json.Unmarshal(response.Body.Bytes(), &contractors)
	if len(contractors) != 3 {
		t.Errorf("Expected contractors retrieved to be 3, found " + strconv.Itoa(len(contractors)))
	}
}

func TestGetContractorCompany(t *testing.T) {
	FreshDatabase()
	addCompanies(2)
	addContractors(1)

	req, _ := http.NewRequest("GET", "/contractor/1/company", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != 1.0 {
		t.Errorf("Expected the id to remain the same (%v). Got %v", 1.0, m["id"])
	}

	if m["name"] != "Company 0" {
		t.Errorf("Expected the name to be '%v'. Got '%v'", "Company 0", m["name"])
	}
}