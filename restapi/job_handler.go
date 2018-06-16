package restapi

import (
	"database/sql"
	"net/http"
	"strconv"
	"upsizeAPI/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"time"
	"fmt"
)

func (a *Api) createJob(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("create job", startTime)
	role := r.Header.Get("authRole")
	if !IsManagerOrAdmin(role) {
		respondWithError(w, http.StatusUnauthorized, "You need to be a manager or admin")
		return
	}
	var j models.Job
	decoder := json.NewDecoder(r.Body)

	if role == "manager" {
		j.ManagerID, _ = strconv.Atoi(r.Header.Get("authId"))
	}

	if err := decoder.Decode(&j); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	if role == "manager" {
		managerId, _ := strconv.Atoi(r.Header.Get("authId"))
		j.ManagerID = managerId
	}

	if err := j.CreateJob(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, j)
}

func (a *Api) updateJob(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("update job", startTime)
	if !IsManagerOrAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin or manager access required")
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Job ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: strconv.Itoa(models.GetCompanyFromJobID(a.DB, vars["id"]))}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	var j models.Job
	if !validPayload(w, r, &j) {
		return
	}
	defer r.Body.Close()
	j.ID = id

	if err := j.UpdateJob(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, j)
}

func (a *Api) deleteJob(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("delete job", startTime)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: strconv.Itoa(models.GetCompanyFromJobID(a.DB, vars["id"]))}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	j := models.Job{ID: id}
	if err := j.DeleteJob(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Api) getJob(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get job", startTime)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}
	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: strconv.Itoa(models.GetCompanyFromJobID(a.DB, vars["id"]))}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	j := models.Job{ID: id}
	if err := j.GetJob(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "job not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, j)
}

func (a *Api) getJobs(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get jobs", startTime)

	role := r.Header.Get("authRole")
	if !IsManagerOrAdmin(role) {
		respondWithError(w, http.StatusUnauthorized, "You need to be a manager or admin")
		return
	}

	companyId := ""
	if role == "manager" {
		companyId = r.Header.Get("authCompanyId")
		fmt.Println(companyId)
	}

	Jobs, err := models.GetJobs(a.DB, companyId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, Jobs)
}

func (a *Api) getJobContractors(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get job contractors", startTime)

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: strconv.Itoa(models.GetCompanyFromJobID(a.DB, mux.Vars(r)["id"]))}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	contractors, err := models.GetJobContractors(a.DB, mux.Vars(r)["id"])

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, contractors)
}
