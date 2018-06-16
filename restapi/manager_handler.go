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

func (a *Api) createManager(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("create manager", startTime)
	if !IsManagerOrAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin or manager access required")
	}
	var m models.Manager
	if !validPayload(w, r, &m) {
		return
	}
	defer r.Body.Close()

	if err := m.CreateManager(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, m)
}

func (a *Api) updateManager(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("update manager", startTime)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid manager ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: strconv.Itoa(models.GetCompanyIDFromID(a.DB, vars["id"], "manager"))}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	var m models.Manager
	if !validPayload(w, r, &m) {
		return
	}

	defer r.Body.Close()
	m.ID = id

	if err := m.UpdateManager(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func (a *Api) deleteManager(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("delete manager", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid manager ID")
		return
	}

	m := models.Manager{ID: id}
	if err := m.DeleteManager(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Api) getManager(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get manager", startTime)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid manager ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: strconv.Itoa(models.GetCompanyIDFromID(a.DB, vars["id"], "manager"))}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	m := models.Manager{ID: id}
	if err := m.GetManager(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Manager not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func (a *Api) getManagers(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get managers", startTime)
	if !IsManagerOrAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin or manager access required")
	}

	companyId := r.Header.Get("authCompanyId")
	if r.Header.Get("authRole") == "admin" {
		companyId = ""
	}
	managers, err := models.GetManagers(a.DB, companyId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, managers)
}

func (a *Api) getManagerCompany(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get manager company", startTime)

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid manager ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: strconv.Itoa(models.GetCompanyIDFromID(a.DB, vars["id"], "manager"))}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	m := models.Company{}
	if err := m.GetManagerCompany(a.DB, vars["id"]); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Manager company not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, m)
}

func (a *Api) getManagerJobs(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get manager jobs", startTime)

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: strconv.Itoa(models.GetCompanyIDFromID(a.DB, mux.Vars(r)["id"], "manager"))}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	managerJobs, err := models.GetManagerJobs(a.DB, mux.Vars(r)["id"], strings.Split(r.FormValue("status"), ","))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, managerJobs)
}