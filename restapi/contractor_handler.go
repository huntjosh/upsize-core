package restapi

import (
	"database/sql"
	"net/http"
	"strconv"
	"upsizeAPI/models"
	"github.com/gorilla/mux"
	"time"
)

func (a *Api) getContractorCompany(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get contractor company", startTime)

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contractor ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "contractor", OwnerID: vars["id"]}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	m := models.Company{}
	if err := m.GetContractorCompany(a.DB, vars["id"]); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Contractor company not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func (a *Api) createContractor(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("create contractor", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	var c models.Contractor
	if !validPayload(w, r, &c) {
		return
	}
	defer r.Body.Close()

	if err := c.CreateContractor(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

func (a *Api) updateContractor(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("update contractor", startTime)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contractor ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "contractor", OwnerID: vars["id"]}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	var c models.Contractor
	if !validPayload(w, r, &c) {
		return
	}

	defer r.Body.Close()
	c.ID = id

	if err := c.UpdateContractor(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *Api) deleteContractor(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("delete contractor", startTime)
	if r.Header.Get("authRole") != "admin" {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contractor ID")
		return
	}

	c := models.Contractor{ID: id}
	if err := c.DeleteContractor(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Api) getContractor(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get contractor", startTime)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contractor ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "contractor", OwnerID: vars["id"]}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	c := models.Contractor{ID: id}
	if err := c.GetContractor(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Contractor not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *Api) getContractors(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get contractors", startTime)

	if !IsManagerOrAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin or manager access required")
	}
	companyId := ""
	if r.Header.Get("authRole") == "manager" {
		companyId = r.Header.Get("authCompanyId")
	}
	contractors, err := models.GetContractors(a.DB, companyId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, contractors)
}