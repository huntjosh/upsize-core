package restapi

import (
	"database/sql"
	"net/http"
	"strconv"
	"upsizeAPI/models"
	"encoding/json"
	"github.com/gorilla/mux"
	"time"
	"strings"
)

func (a *Api) createCompany(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("create company", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	var c models.Company
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := c.CreateCompany(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, c)
}

func (a *Api) updateCompany(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("update company", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	var c models.Company
	if !validPayload(w, r, &c) {
		return
	}

	defer r.Body.Close()
	c.ID = id

	if err := c.UpdateCompany(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *Api) deleteCompany(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("delete company", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	c := models.Company{ID: id}
	if err := c.DeleteCompany(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Api) getCompany(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get company", startTime)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "contractor", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: vars["id"]}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	c := models.Company{ID: id}
	if err := c.GetCompany(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Company not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *Api) getCompanies(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get companies", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}

	companies, err := models.GetCompanies(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, companies)
}

func (a *Api) getCompanySkills(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get company skills", startTime)

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: vars["id"]}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	companies, err := models.GetCompanySkills(a.DB, vars["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, companies)
}

func (a *Api) createCompanySkill(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("create company skill", startTime)

	var cs models.CompanySkill
	if !validPayload(w, r, &cs) {
		return
	}

	defer r.Body.Close()

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: r.Header.Get("authCompanyId")}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	if err := cs.CreateCompanySkill(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, cs)
}

func (a *Api) deleteCompanySkill(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("delete company skill", startTime)
	if !IsAdmin(r.Header.Get("authRole")) {
		respondWithError(w, http.StatusUnauthorized, "Admin access required")
	}
	vars := mux.Vars(r)
	companyId, err := strconv.Atoi(vars["company_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	skillId, err := strconv.Atoi(vars["skill_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid skill ID")
		return
	}

	c := models.CompanySkill{CompanyID: companyId, SkillID: skillId}
	if err := c.DeleteCompanySkill(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *Api) getCompanySkill(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get company skill", startTime)

	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["company_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid company ID")
		return
	}

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: vars["company_id"]}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	_, err = strconv.Atoi(vars["skill_id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid skill ID")
		return
	}
	c := models.Skill{}
	if err := c.GetCompanySkill(a.DB, vars["company_id"], vars["skill_id"]); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "company skill not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}

func (a *Api) getCompanyJobs(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get company jobs", startTime)

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: mux.Vars(r)["id"]}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	companyJobs, err := models.GetCompanyJobs(a.DB, mux.Vars(r)["id"], strings.Split(r.FormValue("status"), ","))

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, companyJobs)
}

func (a *Api) getCompanyContractors(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	defer logFinished("get company contractors", startTime)

	authRole := r.Header.Get("authRole")
	ag := models.AuthGuard{SameUserRole: "", SameCompanyRoles: []string{"manager"}, OverridingRoles: []string{"admin"}}
	authCheck := models.AuthCheck{AccessorRole: authRole, AccessorID: r.Header.Get("authId"),
		OwnerRole: "company", OwnerID: mux.Vars(r)["id"]}

	if !ag.CanAccess(a.DB, authCheck) {
		respondWithError(w, http.StatusUnauthorized, ag.AuthInfo())
		return
	}

	companyContractors, err := models.GetCompanyContractors(a.DB, mux.Vars(r)["id"])

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, companyContractors)
}