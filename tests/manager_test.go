package tests

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"
	"upsizeAPI/models"
)

func TestCreateManager(t *testing.T) {
	FreshDatabase()

	payload := []byte(`
		{"name":"echo","email":"echovvv@gmail.com","phone":"02040490234","company_id":1}`)

	req, _ := http.NewRequest("PUT", "/manager", bytes.NewBuffer(payload))
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "echo" {
		t.Errorf("Expected manager name to be 'echo'. Got '%v'", m["name"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 2.0 {
		t.Errorf("Expected manager ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetNonExistentManager(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/manager/11", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestGetManager(t *testing.T) {
	FreshDatabase()
	addManagers(1, 1)

	req, _ := http.NewRequest("GET", "/manager/1", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addManagers(count, company int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO managers(name, email, phone, company_id) VALUES($1, $2, $3, $4)",
			"Manager "+strconv.Itoa(i), "bob@gmail.com", "020302932", company)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestUpdateManager(t *testing.T) {
	FreshDatabase()
	addManagers(1, 1)

	req, _ := http.NewRequest("GET", "/manager/1", nil)
	response := executeRequest(req, "manager")
	var originalManager map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalManager)

	payload := []byte(`{"name":"echo","email":"echovvv@gmail.com","phone":"020404902134","company_id":2}`)

	req, _ = http.NewRequest("POST", "/manager/1", bytes.NewBuffer(payload))
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	attributes := []string{"name", "email", "phone", "company_id"}
	if m["id" ] != originalManager["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalManager["id"], m["id"])
	}

	for _, attributeName := range attributes {
		if m[attributeName] == originalManager[attributeName] {
			t.Errorf("Expected the %s to change from '%v' to '%v'. Got '%v'", attributeName,
				originalManager[attributeName], m[attributeName], m[attributeName])
		}
	}
}

func TestDeleteManager(t *testing.T) {
	FreshDatabase()
	addManagers(1, 1)

	req, _ := http.NewRequest("GET", "/manager/2", nil)
	response := executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/manager/2", nil)
	response = executeRequest(req, "admin")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/manager/2", nil)
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestGetManagers(t *testing.T) {
	FreshDatabase()
	addManagers(2, 1)

	req, _ := http.NewRequest("GET", "/managers", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var managers []models.Manager
	json.Unmarshal(response.Body.Bytes(), &managers)
	if len(managers) != 3 {
		t.Errorf("Expected managers retrieved to be 3, found " + strconv.Itoa(len(managers)))
	}
}

func TestGetManagerCompany(t *testing.T) {
	FreshDatabase()
	addCompanies(2)
	addManagers(1, 1)

	req, _ := http.NewRequest("GET", "/manager/2/company", nil)
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

func TestGetNoManagerJobs(t *testing.T) {
	FreshDatabase()
	addManagers(1, 1)
	req, _ := http.NewRequest("GET", "/manager/2/jobs", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 0 {
		t.Errorf("Expected jobs retrieved to be 0, found " + strconv.Itoa(len(jobs)))
	}
}

func TestGetManagerJobsNoStatusFilter(t *testing.T) {
	FreshDatabase()
	addManagers(1, 1)
	addJobs(2, "filling", 2)
	req, _ := http.NewRequest("GET", "/manager/2/jobs", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 2 {
		t.Errorf("Expected jobs retrieved to be 2, found " + strconv.Itoa(len(jobs)))
	}
}

func TestGetManagerJobsStatusFilterSingle(t *testing.T) {
	FreshDatabase()
	addManagers(1, 1)
	addJobs(1, "filling", 2)
	addJobs(3, "underway", 2)

	req, _ := http.NewRequest("GET", "/manager/2/jobs?status=filling", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 1 {
		t.Errorf("Expected jobs retrieved to be 1, found " + strconv.Itoa(len(jobs)))
	}
}

func TestGetManagerJobsStatusFilterAll(t *testing.T) {
	FreshDatabase()
	addManagers(1, 1)
	addJobs(1, "filling", 2)
	addJobs(3, "underway", 2)

	req, _ := http.NewRequest("GET", "/manager/2/jobs?status=filling,underway", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 4 {
		t.Errorf("Expected jobs retrieved to be 1, found " + strconv.Itoa(len(jobs)))
	}
}