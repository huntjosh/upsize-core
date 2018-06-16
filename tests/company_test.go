package tests

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"
	"upsizeAPI/models"
)

func TestCreateCompany(t *testing.T) {
	FreshDatabase()

	payload := []byte(`{"name":"best company"}`)

	req, _ := http.NewRequest("PUT", "/company", bytes.NewBuffer(payload))
	response := executeRequest(req, "admin")

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "best company" {
		t.Errorf("Expected company name to be 'best company'. Got '%v'", m["name"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected company ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetWrongCompany(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/company/11", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestGetCompany(t *testing.T) {
	FreshDatabase()
	addCompanies(1)

	req, _ := http.NewRequest("GET", "/company/1", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addCompanies(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO companies(name) VALUES($1)", "Company "+strconv.Itoa(i))
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestUpdateCompany(t *testing.T) {
	FreshDatabase()
	addCompanies(1)

	req, _ := http.NewRequest("GET", "/company/1", nil)
	response := executeRequest(req, "manager")
	var originalCompany map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalCompany)

	payload := []byte(`{"name":"Company - updated name"}`)

	req, _ = http.NewRequest("POST", "/company/1", bytes.NewBuffer(payload))
	response = executeRequest(req, "admin")

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalCompany["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalCompany["id"], m["id"])
	}

	if m["name"] == originalCompany["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalCompany["name"], m["name"], m["name"])
	}
}

func TestDeleteCompany(t *testing.T) {
	FreshDatabase()
	addCompanies(1)

	req, _ := http.NewRequest("GET", "/company/1", nil)
	response := executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/company/1", nil)
	response = executeRequest(req, "admin")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/company/1", nil)
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetCompanies(t *testing.T) {
	FreshDatabase()
	addCompanies(3)

	req, _ := http.NewRequest("GET", "/companies", nil)
	response := executeRequest(req, "admin")

	checkResponseCode(t, http.StatusOK, response.Code)
	var companies []models.Company
	json.Unmarshal(response.Body.Bytes(), &companies)
	if len(companies) != 3 {
		t.Errorf("Expected jobs retrieved to be 2, found " + strconv.Itoa(len(companies)))
	}
}

func TestGetNoCompanyJobs(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/company/1/jobs", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 0 {
		t.Errorf("Expected jobs retrieved to be 0, found " + strconv.Itoa(len(jobs)))
	}
}

func TestGetCompanyJobsNoStatusFilter(t *testing.T) {
	FreshDatabase()
	addManagers(2, 1)
	addJobs(2, "filling", 2)
	addCompanies(2)
	req, _ := http.NewRequest("GET", "/company/1/jobs", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 2 {
		t.Errorf("Expected jobs retrieved to be 2, found " + strconv.Itoa(len(jobs)))
	}
}

func TestGetCompanyJobsStatusFilterSingle(t *testing.T) {
	FreshDatabase()
	addManagers(2, 1)
	addCompanies(2)
	addJobs(1, "filling", 2)
	addJobs(3, "underway", 2)

	req, _ := http.NewRequest("GET", "/company/1/jobs?status=filling", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 1 {
		t.Errorf("Expected jobs retrieved to be 1, found " + strconv.Itoa(len(jobs)))
	}
}

func TestGetCompanyJobsStatusFilterAll(t *testing.T) {
	FreshDatabase()
	addManagers(2, 1)
	addCompanies(2)
	addJobs(1, "filling", 2)
	addJobs(3, "underway", 2)

	req, _ := http.NewRequest("GET", "/company/1/jobs?status=filling,underway", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 4 {
		t.Errorf("Expected jobs retrieved to be 4, found " + strconv.Itoa(len(jobs)))
	}
}

func TestGetCompanyContractors(t *testing.T) {
	FreshDatabase()
	addCompanies(1)
	addContractors(2)

	req, _ := http.NewRequest("GET", "/company/1/contractors", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var contractors []models.Contractor
	json.Unmarshal(response.Body.Bytes(), &contractors)
	if len(contractors) != 3 {
		t.Errorf("Expected jobs retrieved to be 3, found " + strconv.Itoa(len(contractors)))
	}
}

func TestGetCompanyContractorsWrongCompany(t *testing.T) {
	FreshDatabase()
	addCompanies(1)
	addContractors(2)

	req, _ := http.NewRequest("GET", "/company/2/contractors", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}