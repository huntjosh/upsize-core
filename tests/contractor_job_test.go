package tests

import (
	"testing"
	"net/http"
	"bytes"
	"encoding/json"
	"upsizeAPI/models"
	"strconv"
)

func TestCreateContractorJob(t *testing.T) {
	FreshDatabase()

	payload := []byte(`{"contractor_id":1,"status":"invited","state_seen":false,"job_id":1}`)

	req, _ := http.NewRequest("PUT", "/contractor/1/job", bytes.NewBuffer(payload))
	response := executeRequest(req, "admin")

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["contractor_id"] != 1.0 {
		t.Errorf("Expected ContractorJob contractor_id to be '1'. Got '%v'", m["contractor_id"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected ContractorJob id to be '1'. Got '%v'", m["id"])
	}

	if m["job_id"] != 1.0 {
		t.Errorf("Expected ContractorJob job_id to be '1'. Got '%v'", m["job_id"])
	}

	if m["status"] != "invited" {
		t.Errorf("Expected ContractorJob status to be 'invited'. Got '%v'", m["status"])
	}

	if m["state_seen"] != false {
		t.Errorf("Expected ContractorJob state_seen to be 'false'. Got '%v'", m["state_seen"])
	}
}

func TestGetNonExistentContractorJob(t *testing.T) {
	FreshDatabase()

	req, _ := http.NewRequest("GET", "/contractor/1/job/1", nil)
	response := executeRequest(req, "contractor")

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "ContractorJob not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'ContractorJob not found'. Got '%s'", m["error"])
	}
}

func TestGetContractorJob(t *testing.T) {
	FreshDatabase()
	addContractorJobs(1, false)

	req, _ := http.NewRequest("GET", "/contractor/1/job/1", nil)
	response := executeRequest(req, "contractor")

	checkResponseCode(t, http.StatusOK, response.Code)
}

func addContractorJobs(count int, incrementContractorId bool) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		contractorId := 1
		if incrementContractorId {
			contractorId = i + 1
		}
		_, err := a.DB.Exec("INSERT INTO contractor_jobs(contractor_id, status, state_seen, job_id) VALUES($1, $2, $3, $4)",
			contractorId, "invited", false, 1)
		if err != nil {
			panic(err.Error())
		}
	}
}

func addContractorJobUnseenCounts(count int, status string) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		_, err := a.DB.Exec("INSERT INTO contractor_jobs(contractor_id, status, state_seen, job_id) VALUES($1, $2, $3, $4)",
			1, status, false, 1)
		if err != nil {
			panic(err.Error())
		}
	}
}

func TestUpdateContractorJob(t *testing.T) {
	FreshDatabase()
	addContractorJobs(1, false)

	req, _ := http.NewRequest("GET", "/contractor/1/job/1", nil)
	response := executeRequest(req, "manager")
	var originalContractorJob map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalContractorJob)

	payload := []byte(`{"contractor_id":1,"status":"requesting","state_seen":true,"job_id":1}`)

	req, _ = http.NewRequest("POST", "/contractor/1/job/1", bytes.NewBuffer(payload))
	response = executeRequest(req, "contractor")
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	attributes := []string{"status", "state_seen"}

	for _, attributeName := range attributes {
		if m[attributeName] == originalContractorJob[attributeName] {
			t.Errorf("Expected the %s to change from '%v' to '%v'. Got '%v'", attributeName,
				originalContractorJob[attributeName], m[attributeName], m[attributeName])
		}
	}
}

func TestDeleteContractorJob(t *testing.T) {
	FreshDatabase()
	addContractorJobs(1, false)

	req, _ := http.NewRequest("GET", "/contractor/1/job/1", nil)
	response := executeRequest(req, "contractor")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/contractor/1/job/1", nil)
	response = executeRequest(req, "admin")
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/contractor/1/job/1", nil)
	response = executeRequest(req, "contractor")
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetContractorJobs(t *testing.T) {
	FreshDatabase()
	addContractorJobs(3, false)

	req, _ := http.NewRequest("GET", "/contractor/1/jobs", nil)
	response := executeRequest(req, "contractor")

	checkResponseCode(t, http.StatusOK, response.Code)
	var contractorJobs []models.ContractorJob
	json.Unmarshal(response.Body.Bytes(), &contractorJobs)
	if len(contractorJobs) != 3 {
		t.Errorf("Expected contractor jobs retrieved to be 3, found " + strconv.Itoa(len(contractorJobs)))
	}
}

func TestGetContractorJobUnseenCounts(t *testing.T) {
	FreshDatabase()
	addContractorJobUnseenCounts(3, "invited")
	addContractorJobUnseenCounts(3, "requesting")

	req, _ := http.NewRequest("GET", "/contractor/1/jobs/unseenCounts", nil)
	response := executeRequest(req, "contractor")

	checkResponseCode(t, http.StatusOK, response.Code)
	var counts []models.StatusUnseenCount
	json.Unmarshal(response.Body.Bytes(), &counts)
	for _, count := range counts {
		if count.Count != 3 {
			t.Errorf("Expected contractors job count per status to be 3, found " + strconv.Itoa(count.Count))
		}
	}

}