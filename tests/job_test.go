package tests

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"strconv"
	"upsizeAPI/models"
)

func TestCreateJob(t *testing.T) {
	FreshDatabase()

	payload := []byte(`
		{"name":"walk dog","effort":"2 days","start_date":"2018-01-08T04:05:06-01:00","status":"filling","description":"Nice job"}`)

	req, _ := http.NewRequest("PUT", "/job", bytes.NewBuffer(payload))
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "walk dog" {
		t.Errorf("Expected job name to be 'walk dog'. Got '%v'", m["name"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected job ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetNonExistentJob(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/job/11", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestGetJob(t *testing.T) {
	FreshDatabase()
	addJobs(1, "filling", 2)
	addManagers(3, 1)

	req, _ := http.NewRequest("GET", "/job/1", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addJobs(count int, status string, company int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO jobs(name, effort, start_date, status, description, manager_id) VALUES($1, $2, $3, $4, $5, $6)",
			"job "+strconv.Itoa(i), "3 weeks", "2018-02-08T04:05:06-01:00", status, "nice joooob", company)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestUpdateJob(t *testing.T) {
	FreshDatabase()
	addManagers(1, 1)
	addJobs(1, "filling", 2)

	req, _ := http.NewRequest("GET", "/job/1", nil)
	response := executeRequest(req, "manager")
	var originaljob map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originaljob)

	payload := []byte(`{"name":"walk dog","effort":"2 days","start_date":"2018-01-08T04:05:06-01:00","status":"underway","description":"Nice job","manager_id":1}`)

	req, _ = http.NewRequest("POST", "/job/1", bytes.NewBuffer(payload))
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	attributes := []string{"name", "effort", "start_date", "status", "description", "manager_id"}
	if m["id" ] != originaljob["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originaljob["id"], m["id"])
	}

	for _, attributeName := range attributes {
		if m[attributeName] == originaljob[attributeName] {
			t.Errorf("Expected the %s to change from '%v' to '%v'. Got '%v'", attributeName,
				originaljob[attributeName], m[attributeName], m[attributeName])
		}
	}
}

func TestDeleteJob(t *testing.T) {
	FreshDatabase()
	addJobs(1, "filling", 2)
	addManagers(1, 1)

	req, _ := http.NewRequest("GET", "/job/1", nil)
	response := executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/job/1", nil)
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/job/1", nil)
	response = executeRequest(req, "manager")
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestGetNoJobs(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/jobs", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 0 {
		t.Errorf("Expected jobs retrieved to be 0, found " + strconv.Itoa(len(jobs)))
	}
}

func TestGetJobs(t *testing.T) {
	FreshDatabase()
	addJobs(2, "filling", 1)
	addManagers(2, 2)
	req, _ := http.NewRequest("GET", "/jobs", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var jobs []models.Job
	json.Unmarshal(response.Body.Bytes(), &jobs)
	if len(jobs) != 2 {
		t.Errorf("Expected jobs retrieved to be 2, found " + strconv.Itoa(len(jobs)))
	}
}

func TestGetJobContractors(t *testing.T) {
	FreshDatabase()
	addJobs(2, "filling", 1)
	addContractors(2)
	addContractorJobs(2, true)

	req, _ := http.NewRequest("GET", "/job/1/contractors", nil)
	response := executeRequest(req, "manager")

	checkResponseCode(t, http.StatusOK, response.Code)
	var contractors []models.Contractor
	json.Unmarshal(response.Body.Bytes(), &contractors)
	if len(contractors) != 2 {
		t.Errorf("Expected contractors retrieved to be 2, found " + strconv.Itoa(len(contractors)))
	}
}

