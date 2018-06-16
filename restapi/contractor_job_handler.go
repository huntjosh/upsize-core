package restapi

import (
	"database/sql"
	"net/http"
	"strconv"
	"upsizeAPI/models"
	"github.com/gorilla/mux"
	"time"
	"strings"
)

func (a *Api) getContractorJobs(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get contractor jobs", startTime)

	ownerId := mux.Vars(r)["contractor_id"]

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "contractor", OwnerID: ownerId}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	contractorJobs, err := models.GetContractorJobs(a.DB, ownerId, strings.Split(r.FormValue("status"), ","))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, contractorJobs)
}

func (a *Api) createContractorJob(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("create contractor job", startTime)

	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}

	var c models.ContractorJob
	if !validPayload(w, r, &c) {
		return
	}

	defer r.Body.Close()

	if err := c.CreateContractorJob(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

func (a *Api) updateContractorJob(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("update contractor job", startTime)
	vars := mux.Vars(r)
	ownerId := vars["contractor_id"]

	jobId, err := strconv.Atoi(vars["job_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	contractorId, err := strconv.Atoi(vars["contractor_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contractor ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "contractor", OwnerID: ownerId}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	var c models.ContractorJob
	if !validPayload(w, r, &c) {
		return
	}
	defer r.Body.Close()
	c.JobID = jobId
	c.ContractorID = contractorId

	if err := c.UpdateContractorJob(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *Api) deleteContractorJob(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("delete contractor job", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	vars := mux.Vars(r)
	contractorId, err := strconv.Atoi(vars["contractor_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contractor ID")
		return
	}

	jobId, err := strconv.Atoi(vars["job_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	c := models.ContractorJob{ContractorID: contractorId, JobID: jobId}
	if err := c.DeleteContractorJob(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Api) getContractorJob(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get contractor job", startTime)

	vars := mux.Vars(r)
	contractorId, err := strconv.Atoi(vars["contractor_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contractor ID")
		return
	}

	jobId, err := strconv.Atoi(vars["job_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "contractor", OwnerID: vars["contractor_id"]}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	c := models.ContractorJob{ContractorID: contractorId, JobID: jobId}
	if err := c.GetContractorJob(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "ContractorJob not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *Api) getContractorJobUnseenCounts(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get contractor job unseen counts", startTime)
	ownerId := mux.Vars(r)["contractor_id"]
	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "contractor", OwnerID: ownerId}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	counts, err := models.GetContractorUnseenCounts(a.DB, mux.Vars(r)["contractor_id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, counts)
}